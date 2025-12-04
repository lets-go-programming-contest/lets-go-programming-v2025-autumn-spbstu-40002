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
	channels    map[string]chan string
	handlers    []handler
	mu          sync.RWMutex
	bufferSize  int
	channelLock sync.RWMutex
	closed      bool
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
		closed:     false,
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	c.channelLock.Lock()
	defer c.channelLock.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.bufferSize)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) getChannel(name string) (chan string, bool) {
	c.channelLock.RLock()
	defer c.channelLock.RUnlock()

	ch, exists := c.channels[name]
	return ch, exists
}

func (c *Conveyer) RegisterDecorator(fn DecoratorFunc, input, output string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputChan := c.getOrCreateChannel(input)
	outputChan := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, handler{
		handlerType: "decorator",
		fn:          fn,
		inputs:      []string{input},
		outputs:     []string{output},
	})

	c.channelLock.Lock()
	c.channels[input] = inputChan
	c.channels[output] = outputChan
	c.channelLock.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputChans := make([]chan string, len(inputs))
	for i, input := range inputs {
		inputChans[i] = c.getOrCreateChannel(input)
	}
	outputChan := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, handler{
		handlerType: "multiplexer",
		fn:          fn,
		inputs:      inputs,
		outputs:     []string{output},
	})

	c.channelLock.Lock()
	for _, input := range inputs {
		c.channels[input] = c.getOrCreateChannel(input)
	}
	c.channels[output] = outputChan
	c.channelLock.Unlock()
}

func (c *Conveyer) RegisterSeparator(fn SeparatorFunc, input string, outputs []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputChan := c.getOrCreateChannel(input)
	outputChans := make([]chan string, len(outputs))
	for i, output := range outputs {
		outputChans[i] = c.getOrCreateChannel(output)
	}

	c.handlers = append(c.handlers, handler{
		handlerType: "separator",
		fn:          fn,
		inputs:      []string{input},
		outputs:     outputs,
	})

	c.channelLock.Lock()
	c.channels[input] = inputChan
	for _, output := range outputs {
		c.channels[output] = c.getOrCreateChannel(output)
	}
	c.channelLock.Unlock()
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil // или можно вернуть ошибку
	}
	c.mu.Unlock()

	g, ctx := errgroup.WithContext(ctx)

	c.mu.RLock()
	handlers := make([]handler, len(c.handlers))
	copy(handlers, c.handlers)
	c.mu.RUnlock()

	for _, h := range handlers {
		h := h

		switch h.handlerType {
		case "decorator":
			if fn, ok := h.fn.(DecoratorFunc); ok {
				g.Go(func() error {
					inputChan, _ := c.getChannel(h.inputs[0])
					outputChan, _ := c.getChannel(h.outputs[0])
					return fn(ctx, inputChan, outputChan)
				})
			}

		case "multiplexer":
			if fn, ok := h.fn.(MultiplexerFunc); ok {
				g.Go(func() error {
					inputChans := make([]chan string, len(h.inputs))
					for i, input := range h.inputs {
						inputChans[i], _ = c.getChannel(input)
					}
					outputChan, _ := c.getChannel(h.outputs[0])
					return fn(ctx, inputChans, outputChan)
				})
			}

		case "separator":
			if fn, ok := h.fn.(SeparatorFunc); ok {
				g.Go(func() error {
					inputChan, _ := c.getChannel(h.inputs[0])
					outputChans := make([]chan string, len(h.outputs))
					for i, output := range h.outputs {
						outputChans[i], _ = c.getChannel(output)
					}
					return fn(ctx, inputChan, outputChans)
				})
			}
		}
	}

	err := g.Wait()

	c.closeAllChannels()

	return err
}

func (c *Conveyer) closeAllChannels() {
	c.channelLock.Lock()
	defer c.channelLock.Unlock()

	// Помечаем как закрытый
	c.mu.Lock()
	c.closed = true
	c.mu.Unlock()

	for name, ch := range c.channels {
		// Безопасное закрытие канала
		func() {
			defer func() {
				// Восстанавливаемся если канал уже закрыт
				recover()
			}()
			close(ch)
		}()
		delete(c.channels, name)
	}
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return ErrChanNotFound
	}
	c.mu.Unlock()

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
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return "", ErrChanNotFound
	}
	c.mu.Unlock()

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
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return false
	}
	c.mu.Unlock()

	_, exists := c.getChannel(name)
	return exists
}
