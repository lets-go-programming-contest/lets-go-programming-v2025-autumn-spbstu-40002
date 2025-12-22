package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

const Undefined = "undefined"

type Conveyer struct {
	bufferSize int

	mu       sync.RWMutex
	channels map[string]chan string

	wg       sync.WaitGroup
	handlers []func(ctx context.Context) error
}

func New(size int) *Conveyer {
	if size < 0 {
		size = 0
	}

	return &Conveyer{
		bufferSize: size,
		channels:   make(map[string]chan string),
		handlers:   make([]func(ctx context.Context) error, 0),
		mu:         sync.RWMutex{},
		wg:         sync.WaitGroup{},
	}
}

func (c *Conveyer) ensureChannel(id string) chan string {
	c.mu.RLock()
	channel, ok := c.channels[id]
	c.mu.RUnlock()

	if ok {
		return channel
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	channel, ok = c.channels[id]
	if ok {
		return channel
	}

	channel = make(chan string, c.bufferSize)
	c.channels[id] = channel

	return channel
}

func (c *Conveyer) RegisterDecorator(
	handlerFunc func(context.Context, chan string, chan string) error,
	inputID string,
	outputID string,
) {
	inputChan := c.ensureChannel(inputID)
	outputChan := c.ensureChannel(outputID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChan, outputChan)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	handlerFunc func(context.Context, []chan string, chan string) error,
	inputIDs []string,
	outputID string,
) {
	inputChans := make([]chan string, 0, len(inputIDs))
	for _, id := range inputIDs {
		inputChans = append(inputChans, c.ensureChannel(id))
	}

	outputChan := c.ensureChannel(outputID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChans, outputChan)
	})
}

func (c *Conveyer) RegisterSeparator(
	handlerFunc func(context.Context, chan string, []chan string) error,
	inputID string,
	outputIDs []string,
) {
	inputChan := c.ensureChannel(inputID)

	outputChans := make([]chan string, 0, len(outputIDs))
	for _, id := range outputIDs {
		outputChans = append(outputChans, c.ensureChannel(id))
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChan, outputChans)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	for _, registeredHandler := range c.handlers {
		handler := registeredHandler

		c.wg.Add(1)

		go func() {
			defer c.wg.Done()

			_ = handler(ctx)
		}()
	}

	c.wg.Wait()

	c.mu.Lock()
	for _, channel := range c.channels {
		close(channel)
	}
	c.mu.Unlock()

	return nil
}

func (c *Conveyer) Send(inputID string, data string) error {
	c.mu.RLock()
	channel, ok := c.channels[inputID]
	c.mu.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(outputID string) (string, error) {
	c.mu.RLock()
	channel, ok := c.channels[outputID]
	c.mu.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	data, open := <-channel
	if !open {
		return Undefined, nil
	}

	return data, nil
}
