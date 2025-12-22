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
		mu:         sync.RWMutex{},
		channels:   make(map[string]chan string),
		wg:         sync.WaitGroup{},
		handlers:   make([]func(ctx context.Context) error, 0),
	}
}

func (c *Conveyer) ensureChannel(channelID string) chan string {
	c.mu.RLock()
	existingChannel, channelExists := c.channels[channelID]
	c.mu.RUnlock()

	if channelExists {
		return existingChannel
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	existingChannel, channelExists = c.channels[channelID]
	if channelExists {
		return existingChannel
	}

	newChannel := make(chan string, c.bufferSize)
	c.channels[channelID] = newChannel

	return newChannel
}

func (c *Conveyer) RegisterDecorator(
	handlerFunc func(context.Context, chan string, chan string) error,
	inputID string,
	outputID string,
) {
	inputChannel := c.ensureChannel(inputID)
	outputChannel := c.ensureChannel(outputID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannel, outputChannel)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	handlerFunc func(context.Context, []chan string, chan string) error,
	inputIDs []string,
	outputID string,
) {
	inputChannels := make([]chan string, 0, len(inputIDs))
	for _, inputID := range inputIDs {
		inputChannels = append(inputChannels, c.ensureChannel(inputID))
	}

	outputChannel := c.ensureChannel(outputID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannels, outputChannel)
	})
}

func (c *Conveyer) RegisterSeparator(
	handlerFunc func(context.Context, chan string, []chan string) error,
	inputID string,
	outputIDs []string,
) {
	inputChannel := c.ensureChannel(inputID)

	outputChannels := make([]chan string, 0, len(outputIDs))
	for _, outputID := range outputIDs {
		outputChannels = append(outputChannels, c.ensureChannel(outputID))
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannel, outputChannels)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	for _, registeredHandler := range c.handlers {
		currentHandler := registeredHandler

		c.wg.Add(1)
		go func() {
			defer c.wg.Done()

			_ = currentHandler(ctx)
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
	channel, channelExists := c.channels[inputID]
	c.mu.RUnlock()

	if !channelExists {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(outputID string) (string, error) {
	c.mu.RLock()
	channel, channelExists := c.channels[outputID]
	c.mu.RUnlock()

	if !channelExists {
		return "", ErrChanNotFound
	}

	data, open := <-channel
	if !open {
		return Undefined, nil
	}

	return data, nil
}
