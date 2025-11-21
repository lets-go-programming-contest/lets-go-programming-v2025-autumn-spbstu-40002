package conveyer

import (
	"context"
	"fmt"
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

// New creates a new conveyer instance with specified channel buffer size.
func New(size int) *conveyer {
	return &conveyer{
		channels: make(map[string]chan string),
		size:     size,
		handlers: make([]handler, 0),
	}
}

// RegisterDecorator registers a data modifier handler.
func (c *conveyer) RegisterDecorator(
	handlerFn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.ensureChannel(input)
	c.ensureChannel(output)

	c.handlers = append(c.handlers, handler{
		handlerType: "decorator",
		fn:          handlerFn,
		inputs:      []string{input},
		outputs:     []string{output},
	})
}

// RegisterMultiplexer registers a multiplexer handler.
func (c *conveyer) RegisterMultiplexer(
	handlerFn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	for _, input := range inputs {
		c.ensureChannel(input)
	}

	c.ensureChannel(output)

	c.handlers = append(c.handlers, handler{
		handlerType: "multiplexer",
		fn:          handlerFn,
		inputs:      inputs,
		outputs:     []string{output},
	})
}

// RegisterSeparator registers a separator handler.
func (c *conveyer) RegisterSeparator(
	handlerFn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.ensureChannel(input)

	for _, output := range outputs {
		c.ensureChannel(output)
	}

	c.handlers = append(c.handlers, handler{
		handlerType: "separator",
		fn:          handlerFn,
		inputs:      []string{input},
		outputs:     outputs,
	})
}

// Run starts the conveyer and runs all handlers in separate goroutines.
func (c *conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error, len(c.handlers))
	doneChan := make(chan struct{}, len(c.handlers))

	// Start all handlers
	for _, handlerItem := range c.handlers {
		go func(handlerItem handler) {
			err := c.runHandler(ctx, handlerItem)
			if err != nil {
				errChan <- err
			}
			doneChan <- struct{}{}
		}(handlerItem)
	}

	// Wait for all handlers to complete or error/context cancellation
	completed := 0
	for completed < len(c.handlers) {
		select {
		case err := <-errChan:
			if err != nil {
				cancel()
				c.stop()

				return err
			}
		case <-ctx.Done():
			// Stop channels first to allow handlers to finish
			c.stop()
			// Wait for all handlers to finish
			for completed < len(c.handlers) {
				select {
				case <-doneChan:
					completed++
				case <-errChan:
					// Ignore errors when context is cancelled
				}
			}

			return fmt.Errorf("context cancelled: %w", ctx.Err())
		case <-doneChan:
			completed++
		}
	}

	return nil
}

// Send sends data to a channel identified by input ID.
func (c *conveyer) Send(input string, data string) error {
	channel, exists := c.channels[input]
	if !exists {
		return ErrChanNotFound
	}

	select {
	case channel <- data:
		return nil
	default:
		return ErrChanNotFound
	}
}

// Recv receives data from a channel identified by output ID.
func (c *conveyer) Recv(output string) (string, error) {
	channel, exists := c.channels[output]
	if !exists {
		return "", ErrChanNotFound
	}

	data, ok := <-channel
	if !ok {
		return undefined, nil
	}

	return data, nil
}

// ensureChannel creates a channel if it doesn't exist.
func (c *conveyer) ensureChannel(name string) {
	if _, exists := c.channels[name]; !exists {
		c.channels[name] = make(chan string, c.size)
	}
}

// runHandler executes a single handler based on its type.
func (c *conveyer) runHandler(ctx context.Context, handlerItem handler) error {
	switch handlerItem.handlerType {
	case "decorator":
		handlerFn, ok := handlerItem.fn.(func(ctx context.Context, input chan string, output chan string) error)
		if !ok {
			return nil
		}

		inputChan := c.channels[handlerItem.inputs[0]]
		outputChan := c.channels[handlerItem.outputs[0]]

		return handlerFn(ctx, inputChan, outputChan)

	case "multiplexer":
		handlerFn, ok := handlerItem.fn.(func(ctx context.Context, inputs []chan string, output chan string) error)
		if !ok {
			return nil
		}

		inputChans := make([]chan string, len(handlerItem.inputs))

		for i, input := range handlerItem.inputs {
			inputChans[i] = c.channels[input]
		}

		outputChan := c.channels[handlerItem.outputs[0]]

		return handlerFn(ctx, inputChans, outputChan)

	case "separator":
		handlerFn, ok := handlerItem.fn.(func(ctx context.Context, input chan string, outputs []chan string) error)
		if !ok {
			return nil
		}

		inputChan := c.channels[handlerItem.inputs[0]]
		outputChans := make([]chan string, len(handlerItem.outputs))

		for i, output := range handlerItem.outputs {
			outputChans[i] = c.channels[output]
		}

		return handlerFn(ctx, inputChan, outputChans)

	default:
		return nil
	}
}

// stop closes all channels and stops all handlers.
func (c *conveyer) stop() {
	for _, channel := range c.channels {
		close(channel)
	}
}
