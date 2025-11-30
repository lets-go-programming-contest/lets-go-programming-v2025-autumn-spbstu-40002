package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")
var ErrAlreadyClosed = errors.New("already closed")

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

	if ch, exists := c.chans[name]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.chans[name] = ch
	return ch
}

func (c *conveyerImpl) RegisterDecorator(
	handler func(ctx context.Context, input chan string, output chan string) error,
	inputID string,
	outputID string,
) {
	c.ensureChan(inputID)
	c.ensureChan(outputID)

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
	for _, id := range inputIDs {
		c.ensureChan(id)
	}
	c.ensureChan(outputID)

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
	c.ensureChan(inputID)
	for _, id := range outputIDs {
		c.ensureChan(id)
	}

	c.separators = append(
		c.separators,
		separatorReg{fn: handler, inputID: inputID, outputIDs: outputIDs},
	)
}

func (c *conveyerImpl) getChan(name string) chan string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.chans[name]
}

// Разделяем Run на три части
func (c *conveyerImpl) runDecorators(ctx context.Context, launch func(func(context.Context) error)) {
	for _, decorator := range c.decorators {
		inputChan := c.getChan(decorator.inputID)
		outputChan := c.getChan(decorator.outputID)
		fn := decorator.fn
		launch(func(innerCtx context.Context) error {
			select {
			case <-innerCtx.Done():
				return errors.Join(innerCtx.Err())
			default:
				return fn(innerCtx, inputChan, outputChan)
			}
		})
	}
}

func (c *conveyerImpl) runMultiplexers(ctx context.Context, launch func(func(context.Context) error)) {
	for _, mux := range c.multiplexers {
		outputChan := c.getChan(mux.outputID)
		inputChans := make([]chan string, len(mux.inputIDs))
		for i, id := range mux.inputIDs {
			inputChans[i] = c.getChan(id)
		}
		fn := mux.fn
		launch(func(innerCtx context.Context) error {
			select {
			case <-innerCtx.Done():
				return errors.Join(innerCtx.Err())
			default:
				return fn(innerCtx, inputChans, outputChan)
			}
		})
	}
}

func (c *conveyerImpl) runSeparators(ctx context.Context, launch func(func(context.Context) error)) {
	for _, sep := range c.separators {
		inputChan := c.getChan(sep.inputID)
		outputChans := make([]chan string, len(sep.outputIDs))
		for i, id := range sep.outputIDs {
			outputChans[i] = c.getChan(id)
		}
		fn := sep.fn
		launch(func(innerCtx context.Context) error {
			select {
			case <-innerCtx.Done():
				return errors.Join(innerCtx.Err())
			default:
				return fn(innerCtx, inputChan, outputChans)
			}
		})
	}
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

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	launch := func(handler func(context.Context) error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := handler(ctx); err != nil &&
				!errors.Is(err, context.Canceled) &&
				!errors.Is(err, context.DeadlineExceeded) {
				select {
				case errCh <- err:
				default:
				}
			}
		}()
	}

	c.runDecorators(ctx, launch)
	c.runMultiplexers(ctx, launch)
	c.runSeparators(ctx, launch)

	wg.Wait()

	c.mu.Lock()
	if !c.closed {
		for _, ch := range c.chans {
			close(ch)
		}
		c.closed = true
	}
	c.mu.Unlock()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func (c *conveyerImpl) Send(input string, data string) error {
	c.mu.RLock()
	channel, ok := c.chans[input]
	c.mu.RUnlock()
	if !ok {
		return ErrChanNotFound
	}

	channel <- data
	return nil
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	c.mu.RLock()
	channel, ok := c.chans[output]
	c.mu.RUnlock()
	if !ok {
		return "", ErrChanNotFound
	}

	value, ok := <-channel
	if !ok {
		return "undefined", nil
	}

	return value, nil
}
