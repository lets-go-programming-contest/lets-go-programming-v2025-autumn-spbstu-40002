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

	mu sync.RWMutex
}

func New(channelSize int) *Conveyer {
	return &Conveyer{
		channelSize:  channelSize,
		channels:     make(map[string]chan string),
		handlersPool: make([]func(context.Context) error, 0),
	}
}

func (c *Conveyer) makeChannel(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.channels[name]; !ok {
		c.channels[name] = make(chan string, c.channelSize)
	}
}

func (c *Conveyer) makeChannels(names ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, n := range names {
		if _, ok := c.channels[n]; !ok {
			c.channels[n] = make(chan string, c.channelSize)
		}
	}
}

func (c *Conveyer) RegisterDecorator(
	handler func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mu.Lock()
	if _, ok := c.channels[input]; !ok {
		c.channels[input] = make(chan string, c.channelSize)
	}
	if _, ok := c.channels[output]; !ok {
		c.channels[output] = make(chan string, c.channelSize)
	}
	inCh := c.channels[input]
	outCh := c.channels[output]

	c.handlersPool = append(c.handlersPool, func(ctx context.Context) error {
		return handler(ctx, inCh, outCh)
	})
	c.mu.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	handler func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	for _, n := range inputs {
		if _, ok := c.channels[n]; !ok {
			c.channels[n] = make(chan string, c.channelSize)
		}
	}
	if _, ok := c.channels[output]; !ok {
		c.channels[output] = make(chan string, c.channelSize)
	}

	ins := make([]chan string, 0, len(inputs))
	for _, n := range inputs {
		ins = append(ins, c.channels[n])
	}
	out := c.channels[output]

	c.handlersPool = append(c.handlersPool, func(ctx context.Context) error {
		return handler(ctx, ins, out)
	})
	c.mu.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	handler func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	if _, ok := c.channels[input]; !ok {
		c.channels[input] = make(chan string, c.channelSize)
	}
	for _, n := range outputs {
		if _, ok := c.channels[n]; !ok {
			c.channels[n] = make(chan string, c.channelSize)
		}
	}

	in := c.channels[input]
	outs := make([]chan string, 0, len(outputs))
	for _, n := range outputs {
		outs = append(outs, c.channels[n])
	}

	c.handlersPool = append(c.handlersPool, func(ctx context.Context) error {
		return handler(ctx, in, outs)
	})
	c.mu.Unlock()
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.mu.RLock()
	handlers := make([]func(context.Context) error, len(c.handlersPool))
	copy(handlers, c.handlersPool)
	c.mu.RUnlock()

	errGroup, egCtx := errgroup.WithContext(ctx)

	for _, handlerFunc := range handlers {
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
	ch, ok := c.channels[input]
	c.mu.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	defer func() {
		if r := recover(); r != nil {
			err = ErrChanClosed
		}
	}()

	ch <- data
	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, ok := c.channels[output]
	c.mu.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	value, open := <-ch
	if !open {
		return undefinedData, nil
	}

	return value, nil
}
