package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const undefined = "undefined"

var ErrChanNotFound = errors.New("chan not found")

type conveyer struct {
	size     int
	channels map[string]chan string
	handlers []func(ctx context.Context) error
	mu       sync.RWMutex
}

func New(size int) *conveyer {
	return &conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: []func(ctx context.Context) error{},
		mu:       sync.RWMutex{},
	}
}

func (c *conveyer) RegisterDecorator(
	funktion func(ctx context.Context, input, output chan string) error,
	input, output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inCh := c.reservedChannel(input)
	outCh := c.reservedChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return funktion(ctx, inCh, outCh)
	})
}

func (c *conveyer) RegisterMultiplexer(
	funktion func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string, output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	outputChan := c.reservedChannel(output)
	inputChan := make([]chan string, len(inputs))

	for index, name := range inputs {
		inputChan[index] = c.reservedChannel(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return funktion(ctx, inputChan, outputChan)
	})
}

func (c *conveyer) RegisterSeparator(
	funktion func(ctx context.Context, input chan string, outputs []chan string) error,
	input string, outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputChan := c.reservedChannel(input)
	outputChan := make([]chan string, len(outputs))

	for index, name := range outputs {
		outputChan[index] = c.reservedChannel(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return funktion(ctx, inputChan, outputChan)
	})
}

func (c *conveyer) Send(input, data string) error {
	c.mu.RLock()
	chanel, ok := c.channels[input]
	c.mu.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	chanel <- data

	return nil
}

func (c *conveyer) Run(ctx context.Context) error {
	defer c.channelsClose()

	errctx, ctx := errgroup.WithContext(ctx)

	for _, handler := range c.handlers {
		errctx.Go(func() error {
			return handler(ctx)
		})
	}

	if err := errctx.Wait(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (c *conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	chanel, noErr := c.channels[output]
	c.mu.RUnlock()

	if !noErr {
		return "", ErrChanNotFound
	}

	value, noErr := <-chanel
	if !noErr {
		return undefined, nil
	}

	return value, nil
}

func (c *conveyer) reservedChannel(name string) chan string {
	if channel, ok := c.channels[name]; ok {
		return channel
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel

	return channel
}

func (c *conveyer) channelsClose() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, channel := range c.channels {
		close(channel)
	}
}
