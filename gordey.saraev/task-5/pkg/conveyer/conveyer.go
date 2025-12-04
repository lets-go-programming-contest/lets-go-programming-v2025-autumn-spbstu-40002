package conveyer

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

type DecoratorFunc func(ctx context.Context, input chan string, output chan string) error
type MultiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error
type SeparatorFunc func(ctx context.Context, input chan string, outputs []chan string) error

type Conveyer struct {
	bufferSize int
	channels   map[string]chan string
	handlers   []func(ctx context.Context) error
	mu         sync.RWMutex
}

func New(bufferSize int) *Conveyer {
	if bufferSize < 0 {
		bufferSize = 0
	}

	return &Conveyer{
		bufferSize: bufferSize,
		channels:   make(map[string]chan string),
		handlers:   make([]func(ctx context.Context) error, 0),
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.bufferSize)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) getChannel(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, exists := c.channels[name]
	return ch, exists
}

func (c *Conveyer) RegisterDecorator(fn DecoratorFunc, input, output string) {
	inputChan := c.getOrCreateChannel(input)
	outputChan := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChan, outputChan)
	})
}

func (c *Conveyer) RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string) {
	inputChans := make([]chan string, len(inputs))
	for i, input := range inputs {
		inputChans[i] = c.getOrCreateChannel(input)
	}
	outputChan := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChans, outputChan)
	})
}

func (c *Conveyer) RegisterSeparator(fn SeparatorFunc, input string, outputs []string) {
	inputChan := c.getOrCreateChannel(input)
	outputChans := make([]chan string, len(outputs))
	for i, output := range outputs {
		outputChans[i] = c.getOrCreateChannel(output)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChan, outputChans)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	if len(c.handlers) == 0 {
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, handler := range c.handlers {
		handler := handler
		g.Go(func() error {
			return handler(ctx)
		})
	}

	return g.Wait()
}

func (c *Conveyer) Send(input string, data string) error {
	ch, exists := c.getChannel(input)
	if !exists {
		return ErrChanNotFound
	}

	select {
	case ch <- data:
		return nil
	default:
		return ErrChanFull
	}
}

func (c *Conveyer) Recv(output string) (string, error) {
	ch, exists := c.getChannel(output)
	if !exists {
		return "", ErrChanNotFound
	}

	select {
	case data, ok := <-ch:
		if !ok {
			return "undefined", nil
		}
		return data, nil
	default:
		return "", ErrNoData
	}
}

func (c *Conveyer) HasChannel(name string) bool {
	_, exists := c.getChannel(name)
	return exists
}
