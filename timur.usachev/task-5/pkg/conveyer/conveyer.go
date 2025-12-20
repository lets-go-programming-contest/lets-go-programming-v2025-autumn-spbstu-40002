package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const undefinedData = "undefined"

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrChanClosed   = errors.New("chan closed")
)

type Conveyer struct {
	channelSize  int
	channels     map[string]chan string
	handlersPool []func(context.Context) error
	mu           sync.RWMutex
}

func New(channelSize int) *Conveyer {
	return &Conveyer{
		channelSize:  channelSize,
		channels:     make(map[string]chan string),
		handlersPool: make([]func(context.Context) error, 0),
		mu:           sync.RWMutex{},
	}
}

func (c *Conveyer) RegisterDecorator(
	handler func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mu.Lock()

	names := []string{input, output}

	for _, n := range names {
		if _, ok := c.channels[n]; !ok {
			c.channels[n] = make(chan string, c.channelSize)
		}
	}

	inputChannel := c.channels[input]
	outputChannel := c.channels[output]

	c.handlersPool = append(c.handlersPool, func(ctx context.Context) error {
		return handler(ctx, inputChannel, outputChannel)
	})

	c.mu.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	handler func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()

	allNames := make([]string, 0, len(inputs)+1)
	allNames = append(allNames, inputs...)
	allNames = append(allNames, output)

	for _, n := range allNames {
		if _, ok := c.channels[n]; !ok {
			c.channels[n] = make(chan string, c.channelSize)
		}
	}

	inputChannels := make([]chan string, 0, len(inputs))

	for _, n := range inputs {
		inputChannels = append(inputChannels, c.channels[n])
	}

	outputChannel := c.channels[output]

	c.handlersPool = append(c.handlersPool, func(ctx context.Context) error {
		return handler(ctx, inputChannels, outputChannel)
	})

	c.mu.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	handler func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()

	names := make([]string, 0, len(outputs)+1)
	names = append(names, input)
	names = append(names, outputs...)

	for _, n := range names {
		if _, ok := c.channels[n]; !ok {
			c.channels[n] = make(chan string, c.channelSize)
		}
	}

	inputChannel := c.channels[input]

	outputChannels := make([]chan string, 0, len(outputs))

	for _, n := range outputs {
		outputChannels = append(outputChannels, c.channels[n])
	}

	c.handlersPool = append(c.handlersPool, func(ctx context.Context) error {
		return handler(ctx, inputChannel, outputChannels)
	})

	c.mu.Unlock()
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.mu.RLock()

	handlersSnapshot := make([]func(context.Context) error, len(c.handlersPool))
	copy(handlersSnapshot, c.handlersPool)

	c.mu.RUnlock()

	errGroup, egCtx := errgroup.WithContext(ctx)

	for _, handlerFunc := range handlersSnapshot {
		hf := handlerFunc

		errGroup.Go(func() error {
			return hf(egCtx)
		})
	}

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("conveyer handlers: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(input string, data string) (err error) {
	c.mu.RLock()

	channel, ok := c.channels[input]

	c.mu.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	defer func() {
		if r := recover(); r != nil {
			err = ErrChanClosed
		}
	}()

	channel <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()

	channel, ok := c.channels[output]

	c.mu.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	value, open := <-channel

	if !open {
		return undefinedData, nil
	}

	return value, nil
}
