package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

const Undefined = "undefined"

type Conveyer struct {
	size     int
	channels map[string]chan string
	handlers []func(ctx context.Context) error
	mutex    sync.RWMutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: []func(ctx context.Context) error{},
		mutex:    sync.RWMutex{},
	}
}

func (c *Conveyer) reserveChannel(name string) chan string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ch, ok := c.channels[name]
	if !ok {
		ch = make(chan string, c.size)
		c.channels[name] = ch
	}
	return ch
}

func (c *Conveyer) RegisterDecorator(fn func(ctx context.Context, input, output chan string) error, input, output string) {
	inCh := c.reserveChannel(input)
	outCh := c.reserveChannel(output)

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	})
}

func (c *Conveyer) RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string) {
	inChs := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChs[i] = c.reserveChannel(name)
	}
	outCh := c.reserveChannel(output)

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inChs, outCh)
	})
}

func (c *Conveyer) RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string) {
	inCh := c.reserveChannel(input)
	outChs := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChs[i] = c.reserveChannel(name)
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inCh, outChs)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	c.mutex.RLock()
	for _, handler := range c.handlers {
		h := handler // обязательно копируем переменную, чтобы замыкание работало правильно
		group.Go(func() error {
			return h(ctx)
		})
	}
	c.mutex.RUnlock()

	err := group.Wait()

	c.mutex.Lock()
	for _, ch := range c.channels {
		close(ch)
	}
	c.mutex.Unlock()

	return err
}

func (c *Conveyer) Send(input string, data string) error {
	c.mutex.RLock()
	ch, ok := c.channels[input]
	c.mutex.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mutex.RLock()
	ch, ok := c.channels[output]
	c.mutex.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	data, ok := <-ch
	if !ok {
		return Undefined, nil
	}
	return data, nil
}
