package conveyer

import (
	"context"
	"errors"
	"fmt"
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
	fn  func(ctx context.Context, input chan string, output chan string) error
	in  string
	out string
}

type multiplexerReg struct {
	fn  func(ctx context.Context, inputs []chan string, output chan string) error
	ins []string
	out string
}

type separatorReg struct {
	fn   func(ctx context.Context, input chan string, outputs []chan string) error
	in   string
	outs []string
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

func (c *conveyerImpl) getChan(name string) chan string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.chans[name]
}

func (c *conveyerImpl) RegisterDecorator(
	handler func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.ensureChan(input)
	c.ensureChan(output)
	c.decorators = append(c.decorators, decoratorReg{fn: handler, in: input, out: output})
}

func (c *conveyerImpl) RegisterMultiplexer(
	handler func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	for _, name := range inputs {
		c.ensureChan(name)
	}
	c.ensureChan(output)
	c.multiplexers = append(c.multiplexers, multiplexerReg{fn: handler, ins: inputs, out: output})
}

func (c *conveyerImpl) RegisterSeparator(
	handler func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.ensureChan(input)
	for _, name := range outputs {
		c.ensureChan(name)
	}
	c.separators = append(c.separators, separatorReg{fn: handler, in: input, outs: outputs})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return ErrAlreadyClosed
	}
	c.mu.RUnlock()

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var waitGroup sync.WaitGroup
	errCh := make(chan error, 1)

	launch := func(fn func(context.Context) error) {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			if err := fn(runCtx); err != nil &&
				!errors.Is(err, context.Canceled) &&
				!errors.Is(err, context.DeadlineExceeded) {
				select {
				case errCh <- err:
				default:
				}
			}
		}()
	}

	for _, decorator := range c.decorators {
		inputChan := c.getChan(decorator.in)
		outputChan := c.getChan(decorator.out)
		handlerFunc := decorator.fn
		launch(func(ctx context.Context) error {
			return handlerFunc(ctx, inputChan, outputChan)
		})
	}

	for _, multiplexer := range c.multiplexers {
		outputChan := c.getChan(multiplexer.out)
		inputChans := make([]chan string, len(multiplexer.ins))
		for i, name := range multiplexer.ins {
			inputChans[i] = c.getChan(name)
		}
		handlerFunc := multiplexer.fn
		launch(func(ctx context.Context) error {
			return handlerFunc(ctx, inputChans, outputChan)
		})
	}

	for _, separator := range c.separators {
		inputChan := c.getChan(separator.in)
		outputChans := make([]chan string, len(separator.outs))
		for i, name := range separator.outs {
			outputChans[i] = c.getChan(name)
		}
		handlerFunc := separator.fn
		launch(func(ctx context.Context) error {
			return handlerFunc(ctx, inputChan, outputChans)
		})
	}

	waitGroup.Wait()

	c.mu.Lock()
	if !c.closed {
		for _, channel := range c.chans {
			close(channel)
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
