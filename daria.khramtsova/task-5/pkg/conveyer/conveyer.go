package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var errChanNotFound = errors.New("chan not found")

type Conveyer struct {
	channelSize int
	channels    map[string]chan string
	handlers    []func(context.Context) error
	mu          sync.RWMutex
}

func New(size int) *Conveyer {
	if size < 0 {
		size = 0
	}

	return &Conveyer{
		channelSize: size,
		channels:    make(map[string]chan string),
		handlers:    make([]func(context.Context) error, 0),
		mu:          sync.RWMutex{},
	}
}

func (c *Conveyer) getOrCreate(name string) chan string {
	c.mu.RLock()
	channel, exists := c.channels[name]
	c.mu.RUnlock()

	if exists {
		return channel
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	channel, exists = c.channels[name]
	if exists {
		return channel
	}

	channel = make(chan string, c.channelSize)
	c.channels[name] = channel

	return channel
}

func (c *Conveyer) closeChannels() {
	for _, channel := range c.channels {
		close(channel)
	}
}

func (c *Conveyer) RegisterDecorator(
	handler func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	inputChannel := c.getOrCreate(input)
	outputChannel := c.getOrCreate(output)

	c.mu.Lock()
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, inputChannel, outputChannel)
	})
	c.mu.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	handler func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	inputChannels := make([]chan string, len(inputs))
	for index, name := range inputs {
		inputChannels[index] = c.getOrCreate(name)
	}

	outputChannel := c.getOrCreate(output)

	c.mu.Lock()
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, inputChannels, outputChannel)
	})
	c.mu.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	handler func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	inputChannel := c.getOrCreate(input)

	outputChannels := make([]chan string, len(outputs))
	for index, name := range outputs {
		outputChannels[index] = c.getOrCreate(name)
	}

	c.mu.Lock()
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, inputChannel, outputChannels)
	})
	c.mu.Unlock()
}

func (c *Conveyer) Run(ctx context.Context) error {
	group, groupCtx := errgroup.WithContext(ctx)

	c.mu.RLock()
	handlers := append([]func(context.Context) error{}, c.handlers...)
	c.mu.RUnlock()

	for _, handler := range handlers {
		group.Go(func() error {
			return handler(groupCtx)
		})
	}

	err := group.Wait()

	c.mu.Lock()
	c.closeChannels()
	c.mu.Unlock()

	if err != nil {
		return fmt.Errorf("conveyer run failed: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(name string, data string) error {
	c.mu.RLock()
	channel, exists := c.channels[name]
	c.mu.RUnlock()

	if !exists {
		return errChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	c.mu.RLock()
	channel, exists := c.channels[name]
	c.mu.RUnlock()

	if !exists {
		return "", errChanNotFound
	}

	data, open := <-channel
	if !open {
		return "undefined", nil
	}

	return data, nil
}
