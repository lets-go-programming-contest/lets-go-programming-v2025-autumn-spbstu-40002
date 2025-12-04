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
	channels   map[string]chan string
	handlers   []handler
	mu         sync.RWMutex
	bufferSize int
}

type handler struct {
	handlerType string
	fn          interface{}
	inputs      []string
	outputs     []string
}

func New(bufferSize int) *Conveyer {
	return &Conveyer{
		channels:   make(map[string]chan string),
		bufferSize: bufferSize,
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
	c.mu.Lock()
	defer c.mu.Unlock()

	c.getOrCreateChannel(input)
	c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, handler{
		handlerType: "decorator",
		fn:          fn,
		inputs:      []string{input},
		outputs:     []string{output},
	})
}

func (c *Conveyer) RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, input := range inputs {
		c.getOrCreateChannel(input)
	}
	c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, handler{
		handlerType: "multiplexer",
		fn:          fn,
		inputs:      inputs,
		outputs:     []string{output},
	})
}

func (c *Conveyer) RegisterSeparator(fn SeparatorFunc, input string, outputs []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.getOrCreateChannel(input)
	for _, output := range outputs {
		c.getOrCreateChannel(output)
	}

	c.handlers = append(c.handlers, handler{
		handlerType: "separator",
		fn:          fn,
		inputs:      []string{input},
		outputs:     outputs,
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	c.mu.RLock()
	handlers := make([]handler, len(c.handlers))
	copy(handlers, c.handlers)

	// Создаём копию каналов для использования в горутинах
	channelsCopy := make(map[string]chan string)
	for k, v := range c.channels {
		channelsCopy[k] = v
	}
	c.mu.RUnlock()

	for _, h := range handlers {
		h := h

		switch h.handlerType {
		case "decorator":
			if fn, ok := h.fn.(DecoratorFunc); ok {
				g.Go(func() error {
					inputChan := channelsCopy[h.inputs[0]]
					outputChan := channelsCopy[h.outputs[0]]
					return fn(ctx, inputChan, outputChan)
				})
			}

		case "multiplexer":
			if fn, ok := h.fn.(MultiplexerFunc); ok {
				g.Go(func() error {
					inputChans := make([]chan string, len(h.inputs))
					for i, input := range h.inputs {
						inputChans[i] = channelsCopy[input]
					}
					outputChan := channelsCopy[h.outputs[0]]
					return fn(ctx, inputChans, outputChan)
				})
			}

		case "separator":
			if fn, ok := h.fn.(SeparatorFunc); ok {
				g.Go(func() error {
					inputChan := channelsCopy[h.inputs[0]]
					outputChans := make([]chan string, len(h.outputs))
					for i, output := range h.outputs {
						outputChans[i] = channelsCopy[output]
					}
					return fn(ctx, inputChan, outputChans)
				})
			}
		}
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
