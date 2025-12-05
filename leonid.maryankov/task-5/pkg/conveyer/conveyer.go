package conveyer

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

type handlerRunner func(ctx context.Context) error

type channel struct {
	ch chan string
}

type ConveyorImpl struct {
	size int

	mu      sync.RWMutex
	chans   map[string]*channel
	runners []handlerRunner
}

func New(size int) *ConveyorImpl {
	return &ConveyorImpl{
		size:  size,
		chans: make(map[string]*channel),
	}
}

func (c *ConveyorImpl) getOrCreate(id string) chan string {
	if id == "" {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	existing, ok := c.chans[id]
	if !ok {
		newCh := make(chan string, c.size)
		c.chans[id] = &channel{ch: newCh}
		return newCh
	}

	return existing.ch

}

func (c *ConveyorImpl) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	inCh := c.getOrCreate(input)
	outCh := c.getOrCreate(output)

	c.runners = append(c.runners, func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	})

}

func (c *ConveyorImpl) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	inCh := c.getOrCreate(input)

	outChs := make([]chan string, 0, len(outputs))
	for _, id := range outputs {
		outChs = append(outChs, c.getOrCreate(id))
	}

	c.runners = append(c.runners, func(ctx context.Context) error {
		return fn(ctx, inCh, outChs)
	})

}

func (c *ConveyorImpl) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	inChs := make([]chan string, 0, len(inputs))
	for _, id := range inputs {
		inChs = append(inChs, c.getOrCreate(id))
	}

	outCh := c.getOrCreate(output)

	c.runners = append(c.runners, func(ctx context.Context) error {
		return fn(ctx, inChs, outCh)
	})

}

func (c *ConveyorImpl) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for i := range c.runners {
		r := c.runners[i]
		g.Go(func() error {
			return r(ctx)
		})
	}

	err := g.Wait()
	c.closeAll()

	return err

}

func (c *ConveyorImpl) closeAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, chObj := range c.chans {
		func(ch chan string) {
			defer func() { _ = recover() }()
			close(ch)
		}(chObj.ch)
	}

}

var ErrChanNotFound = fmt.Errorf("chan not found")

func (c *ConveyorImpl) Send(input string, data string) error {
	c.mu.RLock()
	chObj, ok := c.chans[input]
	c.mu.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	chObj.ch <- data
	return nil

}

func (c *ConveyorImpl) Recv(output string) (string, error) {
	c.mu.RLock()
	chObj, ok := c.chans[output]
	c.mu.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	value, ok := <-chObj.ch
	if !ok {
		return "undefined", nil
	}

	return value, nil

}
