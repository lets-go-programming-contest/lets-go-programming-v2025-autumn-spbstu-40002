package conveyer

import (
	"context"
	"fmt"
	"sync"
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

type conveyerImpl struct {
	size      int
	channels  map[string]chan string
	handlers  []handler
	mu        sync.RWMutex
	wg        sync.WaitGroup
	errCh     chan error
	isRunning bool
}

type handler struct {
	handlerType string
	fn          interface{}
	inputNames  []string
	outputNames []string
}

func New(size int) Conveyer {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]handler, 0),
		errCh:    make(chan error, 1),
	}
}

func (c *conveyerImpl) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers = append(c.handlers, handler{
		handlerType: "decorator",
		fn:          fn,
		inputNames:  []string{input},
		outputNames: []string{output},
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers = append(c.handlers, handler{
		handlerType: "multiplexer",
		fn:          fn,
		inputNames:  inputs,
		outputNames: []string{output},
	})
}

func (c *conveyerImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers = append(c.handlers, handler{
		handlerType: "separator",
		fn:          fn,
		inputNames:  []string{input},
		outputNames: outputs,
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.isRunning {
		c.mu.Unlock()
		return fmt.Errorf("conveyer already running")
	}
	c.isRunning = true

	for _, h := range c.handlers {
		for _, input := range h.inputNames {
			if _, exists := c.channels[input]; !exists {
				c.channels[input] = make(chan string, c.size)
			}
		}
		for _, output := range h.outputNames {
			if _, exists := c.channels[output]; !exists {
				c.channels[output] = make(chan string, c.size)
			}
		}
	}
	c.mu.Unlock()

	for _, h := range c.handlers {
		c.wg.Add(1)
		go func(h handler) {
			defer c.wg.Done()

			var inputs []chan string
			var outputs []chan string

			c.mu.RLock()
			inputs = make([]chan string, len(h.inputNames))
			for i, name := range h.inputNames {
				inputs[i] = c.channels[name]
			}
			outputs = make([]chan string, len(h.outputNames))
			for i, name := range h.outputNames {
				outputs[i] = c.channels[name]
			}
			c.mu.RUnlock()

			var err error
			switch h.handlerType {
			case "decorator":
				if fn, ok := h.fn.(func(ctx context.Context, input chan string, output chan string) error); ok {
					err = fn(ctx, inputs[0], outputs[0])
				}
			case "multiplexer":
				if fn, ok := h.fn.(func(ctx context.Context, inputs []chan string, output chan string) error); ok {
					err = fn(ctx, inputs, outputs[0])
				}
			case "separator":
				if fn, ok := h.fn.(func(ctx context.Context, input chan string, outputs []chan string) error); ok {
					err = fn(ctx, inputs[0], outputs)
				}
			}

			if err != nil {
				select {
				case c.errCh <- err:
				default:
				}
			}
		}(h)
	}

	go func() {
		c.wg.Wait()
		close(c.errCh)
	}()

	select {
	case <-ctx.Done():
		c.closeAllChannels()
		return ctx.Err()
	case err, ok := <-c.errCh:
		c.closeAllChannels()
		if ok {
			return err
		}
		return nil
	}
}

func (c *conveyerImpl) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, ch := range c.channels {
		close(ch)
		delete(c.channels, name)
	}
	c.isRunning = false
}

func (c *conveyerImpl) Send(input string, data string) error {
	c.mu.RLock()
	ch, exists := c.channels[input]
	isRunning := c.isRunning
	c.mu.RUnlock()

	if !isRunning {
		return fmt.Errorf("conveyer not running")
	}
	if !exists {
		return fmt.Errorf("chan %s not found", input)
	}

	select {
	case ch <- data:
		return nil
	default:
		return fmt.Errorf("channel %s is full", input)
	}
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, exists := c.channels[output]
	isRunning := c.isRunning
	c.mu.RUnlock()

	if !isRunning {
		return "", fmt.Errorf("conveyer not running")
	}
	if !exists {
		return "", fmt.Errorf("chan %s not found", output)
	}

	select {
	case data, ok := <-ch:
		if !ok {
			return "undefined", nil
		}
		return data, nil
	default:
		return "", fmt.Errorf("no data in channel %s", output)
	}
}
