package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

type handlerRunner func(ctx context.Context) error

type channel struct {
	ch chan string
}

type conveyorImpl struct {
	size int

	mu      sync.RWMutex
	chans   map[string]*channel
	runners []handlerRunner
}

func New(size int) *conveyorImpl {
	return &conveyorImpl{
		size:  size,
		chans: make(map[string]*channel),
	}
}

func (c *conveyorImpl) getOrCreate(id string) chan string {
	if id == "" {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	ch, ok := c.chans[id]
	if !ok {
		newCh := make(chan string, c.size)
		c.chans[id] = &channel{ch: newCh}
		return newCh
	}
	return ch.ch
}

func (c *conveyorImpl) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	in := c.getOrCreate(input)
	out := c.getOrCreate(output)

	c.runners = append(c.runners, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *conveyorImpl) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	in := c.getOrCreate(input)

	var outs []chan string
	for _, id := range outputs {
		outs = append(outs, c.getOrCreate(id))
	}

	c.runners = append(c.runners, func(ctx context.Context) error {
		return fn(ctx, in, outs)
	})
}

func (c *conveyorImpl) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	var ins []chan string
	for _, id := range inputs {
		ins = append(ins, c.getOrCreate(id))
	}

	out := c.getOrCreate(output)

	c.runners = append(c.runners, func(ctx context.Context) error {
		return fn(ctx, ins, out)
	})
}

func (c *conveyorImpl) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, r := range c.runners {
		run := r
		g.Go(func() error {
			return run(ctx)
		})
	}

	err := g.Wait()

	c.closeAll()

	return err
}

func (c *conveyorImpl) closeAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ch := range c.chans {
		func(ch chan string) {
			defer func() { recover() }()
			close(ch)
		}(ch.ch)
	}
}

func (c *conveyorImpl) Send(input string, data string) error {
	c.mu.RLock()
	ch, ok := c.chans[input]
	c.mu.RUnlock()
	if !ok {
		return ErrChanNotFound
	}
	ch.ch <- data
	return nil
}

func (c *conveyorImpl) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, ok := c.chans[output]
	c.mu.RUnlock()
	if !ok {
		return "", ErrChanNotFound
	}

	val, ok := <-ch.ch
	if !ok {
		return "undefined", nil
	}
	return val, nil
}
