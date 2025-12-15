package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrNoChannel = errors.New("chan not found")

const undef = "undefined"

type Conveyer struct {
	size     int
	channels map[string]chan string
	handlers []func(context.Context) error
	mu       sync.RWMutex
}

func New(size int) Conveyer {
	return Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: []func(context.Context) error{},
		mu:       sync.RWMutex{},
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	defer func() {
		c.mu.Lock()
		for _, channel := range c.channels {
			close(channel)
		}
		c.mu.Unlock()
	}()

	group, ctx := errgroup.WithContext(ctx)

	for _, handler := range c.handlers {
		group.Go(func() error {
			return handler(ctx)
		})
	}

	err := group.Wait()
	if err != nil {
		return fmt.Errorf("conveyer failes due to: %w", err)
	}

	return nil
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

	data, isOpen := <-channel
	if !isOpen {
		return undef, nil
	}

	return data, nil
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
	c.mu.Lock()
	defer c.mu.Unlock()

	c.makeChannels(input)
	c.makeChannels(output)
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
	c.mu.Lock()
	defer c.mu.Unlock()

	c.makeChannels(inputs...)
	c.makeChannels(output)
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		inputsCh := make([]chan string, 0, len(inputs))

		for _, channel := range inputs {
			inputsCh = append(inputsCh, c.channels[channel])
		}

		return handler(ctx, inputsCh, c.channels[output])
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
	c.mu.Lock()
	defer c.mu.Unlock()

	c.makeChannels(input)
	c.makeChannels(outputs...)
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		outputsCh := make([]chan string, 0, len(outputs))

		for _, channel := range outputs {
			outputsCh = append(outputsCh, c.channels[channel])
		}

		return handler(ctx, c.channels[input], outputsCh)
	})
}
