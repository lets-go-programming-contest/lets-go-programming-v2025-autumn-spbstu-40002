package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrChanNotFound  = errors.New("chan not found")
	ErrAlreadyClosed = errors.New("already closed")
)

type Conveyer interface {
	RegisterDecorator(
		fn func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		fn func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		fn func(ctx context.Context, input chan string, outputs []chan string) error,
		input string,
		outputs []string,
	)

	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type decoratorReg struct {
	fn       func(ctx context.Context, input chan string, output chan string) error
	inputID  string
	outputID string
}

type multiplexerReg struct {
	fn       func(ctx context.Context, inputs []chan string, output chan string) error
	inputIDs []string
	outputID string
}

type separatorReg struct {
	fn        func(ctx context.Context, input chan string, outputs []chan string) error
	inputID   string
	outputIDs []string
}

type conveyerImpl struct {
	size         int
	mu           sync.RWMutex
	chans        map[string]chan string
	decorators   []decoratorReg
	multiplexers []multiplexerReg
	separators   []separatorReg
	closed       bool
}

func New(size int) Conveyer {
	return &conveyerImpl{
		size:         size,
		mu:           sync.RWMutex{},
		chans:        make(map[string]chan string),
		decorators:   nil,
		multiplexers: nil,
		separators:   nil,
		closed:       false,
	}
}

func (c *conveyerImpl) ensureChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, exists := c.chans[name]; exists {
		return channel
	}

	channel := make(chan string, c.size)
	c.chans[name] = channel

	return channel
}

func (c *conveyerImpl) RegisterDecorator(
	handler func(ctx context.Context, input chan string, output chan string) error,
	inputID string,
	outputID string,
) {
	_ = c.ensureChan(inputID)
	_ = c.ensureChan(outputID)

	c.decorators = append(
		c.decorators,
		decoratorReg{fn: handler, inputID: inputID, outputID: outputID},
	)
}

func (c *conveyerImpl) RegisterMultiplexer(
	handler func(ctx context.Context, inputs []chan string, output chan string) error,
	inputIDs []string,
	outputID string,
) {
	for _, inputID := range inputIDs {
		_ = c.ensureChan(inputID)
	}
	_ = c.ensureChan(outputID)

	c.multiplexers = append(
		c.multiplexers,
		multiplexerReg{fn: handler, inputIDs: inputIDs, outputID: outputID},
	)
}

func (c *conveyerImpl) RegisterSeparator(
	handler func(ctx context.Context, input chan string, outputs []chan string) error,
	inputID string,
	outputIDs []string,
) {
	_ = c.ensureChan(inputID)

	for _, outputID := range outputIDs {
		_ = c.ensureChan(outputID)
	}

	c.separators = append(
		c.separators,
		separatorReg{fn: handler, inputID: inputID, outputIDs: outputIDs},
	)
}

func (c *conveyerImpl) runDecorators(
	launch func(func(context.Context) error),
) {
	for _, decorator := range c.decorators {
		inputChan := c.getChan(decorator.inputID)
		outputChan := c.getChan(decorator.outputID)
		handlerFunc := decorator.fn

		launch(func(innerCtx context.Context) error {
			select {
			case <-innerCtx.Done():
				return errors.Join(innerCtx.Err())
			default:
				return handlerFunc(innerCtx, inputChan, outputChan)
			}
		})
	}
}

func (c *conveyerImpl) runMultiplexers(
	launch func(func(context.Context) error),
) {
	for _, mux := range c.multiplexers {
		outputChan := c.getChan(mux.outputID)
		inputChans := make([]chan string, len(mux.inputIDs))

		for idx, id := range mux.inputIDs {
			inputChans[idx] = c.getChan(id)
		}

		handlerFunc := mux.fn

		launch(func(innerCtx context.Context) error {
			select {
			case <-innerCtx.Done():
				return errors.Join(innerCtx.Err())
			default:
				return handlerFunc(innerCtx, inputChans, outputChan)
			}
		})
	}
}

func (c *conveyerImpl) runSeparators(
	launch func(func(context.Context) error),
) {
	for _, sep := range c.separators {
		inputChan := c.getChan(sep.inputID)
		outputChans := make([]chan string, len(sep.outputIDs))

		for idx, id := range sep.outputIDs {
			outputChans[idx] = c.getChan(id)
		}

		handlerFunc := sep.fn
		launch(func(innerCtx context.Context) error {
			select {
			case <-innerCtx.Done():
				return errors.Join(innerCtx.Err())
			default:
				return handlerFunc(innerCtx, inputChan, outputChans)
			}
		})
	}
}

func (c *conveyerImpl) getChan(name string) chan string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.chans[name]
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()

		return ErrAlreadyClosed
	}
	c.mu.RUnlock()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var waitGroup sync.WaitGroup
	errorChannel := make(chan error, 1)

	launch := func(handler func(context.Context) error) {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			if err := handler(ctx); err != nil &&
				!errors.Is(err, context.Canceled) &&
				!errors.Is(err, context.DeadlineExceeded) {
				select {
				case errorChannel <- err:
				default:
				}
			}
		}()
	}

	c.runDecorators(launch)
	c.runMultiplexers(launch)
	c.runSeparators(launch)

	waitGroup.Wait()

	c.mu.Lock()
	if !c.closed {
		for _, channel := range c.chans {
			close(channel)
		}
	}
	c.closed = true
	c.mu.Unlock()

	select {
	case err := <-errorChannel:
		return err
	default:
		return nil
	}
}

func (c *conveyerImpl) Send(input string, data string) error {
	c.mu.RLock()
	channel, exists := c.chans[input]
	c.mu.RUnlock()

	if !exists {
		return ErrChanNotFound
	}

	channel <- data
	return nil
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	c.mu.RLock()
	channel, exists := c.chans[output]
	c.mu.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	value, ok := <-channel
	if !ok {
		return "undefined", nil
	}

	return value, nil
}
