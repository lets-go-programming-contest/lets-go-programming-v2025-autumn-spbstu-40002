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

type conveyerImpl struct {
	size     int
	channels map[string]chan string
	handlers []func(ctx context.Context) error
	mu       sync.RWMutex
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
		handlers: []func(ctx context.Context) error{},
	}
}

func (c *conveyerImpl) RegisterDecorator(
	fn func(ctx context.Context, input, output chan string) error,
	input, output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inCh := c.getOrCreateChannel(input)
	outCh := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string, output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	outCh := c.getOrCreateChannel(output)
	inChs := make([]chan string, len(inputs))

	for i, name := range inputs {
		inChs[i] = c.getOrCreateChannel(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inChs, outCh)
	})
}

func (c *conveyerImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string, outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inCh := c.getOrCreateChannel(input)
	outChs := make([]chan string, len(outputs))

	for i, name := range outputs {
		outChs[i] = c.getOrCreateChannel(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inCh, outChs)
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	defer c.closeAllChannels()

	errgr, ctx := errgroup.WithContext(ctx)

	for _, h := range c.handlers {
		errgr.Go(func() error {
			return h(ctx)
		})
	}

	if err := errgr.Wait(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (c *conveyerImpl) Send(input, data string) error {
	c.mu.RLock()
	ch, ok := c.channels[input]
	c.mu.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	ch <- data

	return nil
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, ok := c.channels[output]
	c.mu.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	val, ok := <-ch
	if !ok {
		return undefined, nil
	}
	return val, nil
}

func (c *conveyerImpl) getOrCreateChannel(name string) chan string {
	if channel, ok := c.channels[name]; ok {
		return channel
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel

	return channel
}

func (c *conveyerImpl) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ch := range c.channels {
		close(ch)
	}
}
