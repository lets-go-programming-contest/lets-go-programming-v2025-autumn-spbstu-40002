package conveyer

import (
	"context"
	"errors"
)

const undefined = "undefined"

type conveyer struct {
	channels map[string]chan string
	size     int
	handlers []handler
}

type handler struct {
	handlerType string
	fn          interface{}
	inputs      []string
	outputs     []string
}

// New creates a new conveyer instance with specified channel buffer size
func New(size int) *conveyer {
	return &conveyer{
		channels: make(map[string]chan string),
		size:     size,
		handlers: make([]handler, 0),
	}
}

// RegisterDecorator registers a data modifier handler
func (c *conveyer) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.ensureChannel(input)
	c.ensureChannel(output)

	c.handlers = append(c.handlers, handler{
		handlerType: "decorator",
		fn:          fn,
		inputs:      []string{input},
		outputs:     []string{output},
	})
}

// RegisterMultiplexer registers a multiplexer handler
func (c *conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	for _, input := range inputs {
		c.ensureChannel(input)
	}
	c.ensureChannel(output)

	c.handlers = append(c.handlers, handler{
		handlerType: "multiplexer",
		fn:          fn,
		inputs:      inputs,
		outputs:     []string{output},
	})
}

// RegisterSeparator registers a separator handler
func (c *conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.ensureChannel(input)
	for _, output := range outputs {
		c.ensureChannel(output)
	}

	c.handlers = append(c.handlers, handler{
		handlerType: "separator",
		fn:          fn,
		inputs:      []string{input},
		outputs:     outputs,
	})
}

// Run starts the conveyer and runs all handlers in separate goroutines
func (c *conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error, len(c.handlers))
	doneChan := make(chan struct{}, len(c.handlers))

	// Start all handlers
	for _, h := range c.handlers {
		go func(h handler) {
			var err error

			switch h.handlerType {
			case "decorator":
				fn := h.fn.(func(ctx context.Context, input chan string, output chan string) error)
				inputChan := c.channels[h.inputs[0]]
				outputChan := c.channels[h.outputs[0]]
				err = fn(ctx, inputChan, outputChan)

			case "multiplexer":
				fn := h.fn.(func(ctx context.Context, inputs []chan string, output chan string) error)
				inputChans := make([]chan string, len(h.inputs))
				for i, input := range h.inputs {
					inputChans[i] = c.channels[input]
				}
				outputChan := c.channels[h.outputs[0]]
				err = fn(ctx, inputChans, outputChan)

			case "separator":
				fn := h.fn.(func(ctx context.Context, input chan string, outputs []chan string) error)
				inputChan := c.channels[h.inputs[0]]
				outputChans := make([]chan string, len(h.outputs))
				for i, output := range h.outputs {
					outputChans[i] = c.channels[output]
				}
				err = fn(ctx, inputChan, outputChans)
			}

			if err != nil {
				errChan <- err
			}
			doneChan <- struct{}{}
		}(h)
	}

	// Wait for all handlers to complete or error/context cancellation
	completed := 0
	for completed < len(c.handlers) {
		select {
		case <-ctx.Done():
			c.stop()
			return ctx.Err()
		case err := <-errChan:
			if err != nil {
				c.stop()
				return err
			}
		case <-doneChan:
			completed++
		}
	}

	return nil
}

// Send sends data to a channel identified by input ID
func (c *conveyer) Send(input string, data string) error {
	ch, exists := c.channels[input]
	if !exists {
		return errors.New("chan not found")
	}

	select {
	case ch <- data:
		return nil
	default:
		return errors.New("chan not found")
	}
}

// Recv receives data from a channel identified by output ID
func (c *conveyer) Recv(output string) (string, error) {
	ch, exists := c.channels[output]
	if !exists {
		return "", errors.New("chan not found")
	}

	data, ok := <-ch
	if !ok {
		return undefined, nil
	}

	return data, nil
}

// ensureChannel creates a channel if it doesn't exist
func (c *conveyer) ensureChannel(name string) {
	if _, exists := c.channels[name]; !exists {
		c.channels[name] = make(chan string, c.size)
	}
}

// stop closes all channels and stops all handlers
func (c *conveyer) stop() {
	for _, ch := range c.channels {
		close(ch)
	}
}
