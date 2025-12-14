package conveyer

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Conveyer struct {
	size     int
	channels map[string]chan string
	handlers []func(ctx context.Context) error
	mu       sync.RWMutex
}

func (c *Conveyer) createChannel(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.channels[name]; !exists {
		c.channels[name] = make(chan string, c.size)
	}
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]func(ctx context.Context) error, 0),
		mu:       sync.RWMutex{},
	}
}

func (c *Conveyer) RegisterDecorator(
	handler func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	c.createChannel(input)
	c.createChannel(output)
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, c.channels[input], c.channels[output])
	})
}

func (c *Conveyer) RegisterMultiplexer(
	handler func(
		ctx context.Context,
		inputs []chan string,
		output chan string,
	) error,
	inputs []string,
	output string,
) {
	for _, input := range inputs {
		c.createChannel(input)
	}

	c.createChannel(output)
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		requiredChannels := make([]chan string, 0, len(inputs))
		for _, input := range inputs {
			requiredChannels = append(requiredChannels, c.channels[input])
		}

		return handler(ctx, requiredChannels, c.channels[output])
	})
}

func (c *Conveyer) RegisterSeparator(
	handler func(
		ctx context.Context,
		input chan string,
		outputs []chan string,
	) error,
	input string,
	outputs []string,
) {
	c.createChannel(input)

	for _, output := range outputs {
		c.createChannel(output)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		requiredChannels := make([]chan string, 0, len(outputs))

		for _, out := range outputs {
			requiredChannels = append(requiredChannels, c.channels[out])
		}

		return handler(ctx, c.channels[input], requiredChannels)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	groupHandlers, gctx := errgroup.WithContext(ctx)

	for _, handler := range c.handlers {
		h := handler

		groupHandlers.Go(func() error {
			return h(gctx)
		})
	}

	return groupHandlers.Wait() //nolint:wrapcheck
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.RLock()
	channel, exists := c.channels[input]
	c.mu.RUnlock()

	if !exists {
		return ErrNoChannel
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	channel, exists := c.channels[output]
	c.mu.RUnlock()

	if !exists {
		return "", ErrNoChannel
	}

	data, ok := <-channel
	if !ok {
		return "undefined", nil
	}

	return data, nil
}
