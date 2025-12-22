package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	channels map[string]chan string
	handlers []func(context.Context) error
	mu       sync.RWMutex
	wg       sync.WaitGroup
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		handlers: []func(context.Context) error{},
	}
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	in, out string,
) {
	c.getOrCreateChannel(in)
	c.getOrCreateChannel(out)
	
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, c.channels[in], c.channels[out])
	})
}

func (c *Conveyer) getOrCreateChannel(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.channels[name]; !ok {
		c.channels[name] = make(chan string, 10)
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	for _, h := range c.handlers {
		c.wg.Add(1)
		go func(handler func(context.Context) error) {
			defer c.wg.Done()
			handler(ctx)
		}(h)
	}
	
	c.wg.Wait()
	c.mu.Lock()
	for _, ch := range c.channels {
		close(ch)
	}
	c.mu.Unlock()
	
	return nil
}

func (c *Conveyer) Send(id string, data string) error {
	c.mu.RLock()
	ch, ok := c.channels[id]
	c.mu.RUnlock()
	
	if !ok {
		return ErrChanNotFound
	}
	ch <- data
	return nil
}

func (c *Conveyer) Recv(id string) (string, error) {
	c.mu.RLock()
	ch, ok := c.channels[id]
	c.mu.RUnlock()
	
	if !ok {
		return "", ErrChanNotFound
	}
	data := <-ch
	return data, nil
}
