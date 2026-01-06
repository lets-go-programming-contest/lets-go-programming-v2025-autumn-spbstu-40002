package conveyer

import (
	"context"
	"errors"
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
	ch, ok := c.channels[name]
	c.mu.RUnlock()

	if ok {
		return ch
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok = c.channels[name]; ok {
		return ch
	}

	ch = make(chan string, c.channelSize)
	c.channels[name] = ch
	
	return ch
}

func (c *Conveyer) closeChannels() {
	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	in := c.getOrCreate(input)
	out := c.getOrCreate(output)

	c.mu.Lock()
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
	c.mu.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	inChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChans[i] = c.getOrCreate(name)
	}

	out := c.getOrCreate(output)

	c.mu.Lock()
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inChans, out)
	})
	c.mu.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	in := c.getOrCreate(input)

	outChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChans[i] = c.getOrCreate(name)
	}

	c.mu.Lock()
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, outChans)
	})
	c.mu.Unlock()
}

func (c *Conveyer) Run(ctx context.Context) error {
	group, gctx := errgroup.WithContext(ctx)

	c.mu.RLock()
	handlers := append([]func(context.Context) error{}, c.handlers...)
	c.mu.RUnlock()

	for _, h := range handlers {
		group.Go(func() error {
			return h(gctx)
		})
	}

	err := group.Wait()

	c.mu.Lock()
	c.closeChannels()
	c.mu.Unlock()

	return err
}

func (c *Conveyer) Send(name string, data string) error {
	c.mu.RLock()
	ch, ok := c.channels[name]
	c.mu.RUnlock()

	if !ok {
		return errChanNotFound
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	c.mu.RLock()
	ch, ok := c.channels[name]
	c.mu.RUnlock()

	if !ok {
		return "", errChanNotFound
	}

	data, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return data, nil
}
