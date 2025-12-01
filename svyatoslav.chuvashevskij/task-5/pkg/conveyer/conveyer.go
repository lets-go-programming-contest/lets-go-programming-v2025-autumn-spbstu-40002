package conveyer

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

const undefined = "undefined"

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	size     int
	channels map[string]chan string
	handlers []func(ctx context.Context) error
}

func New(size int) Conveyer {
	return Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: []func(ctx context.Context) error{},
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	defer func() {
		for _, channel := range c.channels {
			close(channel)
		}
	}()

	group, ctx := errgroup.WithContext(ctx)

	for _, handler := range c.handlers {
		group.Go(func() error { return handler(ctx) })
	}

	return group.Wait()
}

func (c *Conveyer) Send(input string, data string) error {
	channel, exists := c.channels[input]
	if !exists {
		return fmt.Errorf("%w, trying to send to \"%s\" chan", ErrChanNotFound, input)
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	channel, exists := c.channels[output]
	if !exists {
		return "", fmt.Errorf("%w, trying to receive from \"%s\" chan", ErrChanNotFound, output)
	}

	data, isOpen := <-channel
	if !isOpen {
		return undefined, nil
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
	_, exist := c.channels[input]
	if !exist {
		c.channels[input] = make(chan string, c.size)
	}
	_, exist = c.channels[output]
	if !exist {
		c.channels[output] = make(chan string, c.size)
	}

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
	inputChannels := make([]chan string, c.size)
	for i, input := range inputs {
		_, exists := c.channels[input]
		if !exists {
			c.channels[input] = make(chan string, c.size)
		}
		inputChannels[i] = c.channels[input]
	}

	_, exist := c.channels[output]
	if !exist {
		c.channels[output] = make(chan string, c.size)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, inputChannels, c.channels[output])
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
	_, exist := c.channels[input]
	if !exist {
		c.channels[input] = make(chan string, c.size)
	}

	outputChannels := make([]chan string, c.size)
	for i, output := range outputs {
		_, exists := c.channels[output]
		if !exists {
			c.channels[output] = make(chan string, c.size)
		}
		outputChannels[i] = c.channels[output]
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, c.channels[input], outputChannels)
	})
}
