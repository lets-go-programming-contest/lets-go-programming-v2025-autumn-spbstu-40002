package conveyer

import (
	"context"
	"fmt"
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
		mu:         sync.RWMutex{},
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, exists := c.channels[name]; exists {
		return channel
	}

	channel := make(chan string, c.bufferSize)
	c.channels[name] = channel

	return channel
}

func (c *Conveyer) getChannel(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	channel, exists := c.channels[name]

	return channel, exists
}

func (c *Conveyer) RegisterDecorator(callback DecoratorFunc, input, output string) {
	inputChan := c.getOrCreateChannel(input)
	outputChan := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return callback(ctx, inputChan, outputChan)
	})
}

func (c *Conveyer) RegisterMultiplexer(callback MultiplexerFunc, inputs []string, output string) {
	inputChannels := make([]chan string, len(inputs))

	for i, inputName := range inputs {
		inputChannels[i] = c.getOrCreateChannel(inputName)
	}

	outputChannel := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return callback(ctx, inputChannels, outputChannel)
	})
}

func (c *Conveyer) RegisterSeparator(callback SeparatorFunc, input string, outputs []string) {
	inputChannel := c.getOrCreateChannel(input)
	outputChannels := make([]chan string, len(outputs))

	for i, outputName := range outputs {
		outputChannels[i] = c.getOrCreateChannel(outputName)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return callback(ctx, inputChannel, outputChannels)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	if len(c.handlers) == 0 {
		return nil
	}

	errGroup, ctx := errgroup.WithContext(ctx)

	for _, handler := range c.handlers {
		errGroup.Go(func() error {
			return handler(ctx)
		})
	}

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("conveyer run: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(input string, data string) error {
	channel, exists := c.getChannel(input)

	if !exists {
		return ErrChanNotFound
	}

	select {
	case channel <- data:
		return nil
	default:
		return ErrChanFull
	}
}

func (c *Conveyer) Recv(output string) (string, error) {
	channel, exists := c.getChannel(output)

	if !exists {
		return "", ErrChanNotFound
	}

	select {
	case data, ok := <-channel:
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
