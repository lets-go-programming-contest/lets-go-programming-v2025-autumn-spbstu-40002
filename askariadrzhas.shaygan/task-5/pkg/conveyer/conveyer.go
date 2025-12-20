package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var errChannelNotFound = errors.New("chan not found")

const undefined = "undefined"

type ConveyerInterface interface {
	RegisterDecorator(
		function func(
			ctx context.Context,
			input chan string,
			output chan string,
		) error,
		input string,
		output string,
	)

	RegisterMultiplexer(
		function func(
			ctx context.Context,
			inputs []chan string,
			output chan string,
		) error,
		inputs []string,
		output string,
	)

	RegisterSeparator(
		function func(
			ctx context.Context,
			input chan string,
			outputs []chan string,
		) error,
		input string,
		outputs []string,
	)

	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type Conveyer struct {
	channels   map[string]chan string
	bufferSize int
	handlers   []func(ctx context.Context) error
	mu         sync.RWMutex
}

func New(size int) Conveyer {
	return Conveyer{
		channels:   make(map[string]chan string),
		bufferSize: size,
		handlers:   make([]func(ctx context.Context) error, 0),
		mu:         sync.RWMutex{},
	}
}

func (c *Conveyer) ensureChannel(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.channels[name]; !ok {
		c.channels[name] = make(chan string, c.bufferSize)
	}
}

func (c *Conveyer) getChannel(name string) (chan string, error) {
	if ch, ok := c.channels[name]; ok {
		return ch, nil
	}

	return nil, errChannelNotFound
}

func (c *Conveyer) closeChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *Conveyer) addHandler(fn func(ctx context.Context) error) {
	c.handlers = append(c.handlers, fn)
}

func (c *Conveyer) RegisterDecorator(
	function func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string, output string,
) {
	c.ensureChannel(input)
	c.ensureChannel(output)

	c.addHandler(func(ctx context.Context) error {
		c.mu.RLock()
		defer c.mu.RUnlock()

		inCh, _ := c.getChannel(input)
		outCh, _ := c.getChannel(output)

		return function(ctx, inCh, outCh)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	function func(
		ctx context.Context,
		inputs []chan string,
		output chan string,
	) error,
	inputs []string, output string,
) {
	for _, inp := range inputs {
		c.ensureChannel(inp)
	}

	c.ensureChannel(output)

	c.addHandler(func(ctx context.Context) error {
		c.mu.RLock()
		defer c.mu.RUnlock()

		inChannels := make([]chan string, len(inputs))

		for i, inp := range inputs {
			ch, _ := c.getChannel(inp)
			inChannels[i] = ch
		}

		outCh, _ := c.getChannel(output)

		return function(ctx, inChannels, outCh)
	})
}

func (c *Conveyer) RegisterSeparator(
	function func(
		ctx context.Context,
		input chan string,
		outputs []chan string,
	) error,
	input string, outputs []string,
) {
	c.ensureChannel(input)

	for _, out := range outputs {
		c.ensureChannel(out)
	}

	c.addHandler(func(ctx context.Context) error {
		c.mu.RLock()
		defer c.mu.RUnlock()

		inCh, _ := c.getChannel(input)

		outChannels := make([]chan string, len(outputs))

		for i, out := range outputs {
			ch, _ := c.getChannel(out)
			outChannels[i] = ch
		}

		return function(ctx, inCh, outChannels)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	defer c.closeChannels()

	group, ctxWithCancel := errgroup.WithContext(ctx)

	c.mu.RLock()

	for _, handler := range c.handlers {
		h := handler

		group.Go(func() error {
			return h(ctxWithCancel)
		})
	}

	c.mu.RUnlock()

	if err := group.Wait(); err != nil {
		return fmt.Errorf("run function failed: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.RLock()

	inCh, err := c.getChannel(input)
	if err != nil {
		c.mu.RUnlock()

		return err
	}

	c.mu.RUnlock()

	inCh <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()

	outCh, err := c.getChannel(output)
	if err != nil {
		c.mu.RUnlock()

		return "", err
	}

	c.mu.RUnlock()

	if data, ok := <-outCh; ok {
		return data, nil
	} else {
		return undefined, nil
	}
}
