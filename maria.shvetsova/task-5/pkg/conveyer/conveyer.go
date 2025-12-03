package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	errChanNotFound = errors.New("chan not found")
	errNoHandlers   = errors.New("conveyer has no handlers")
)

type Conveyer struct {
	channelSize  int
	channels     map[string]chan string
	handlersPool []func(context.Context) error
	mu           sync.RWMutex
}

func New(channelSize int) *Conveyer {
	return &Conveyer{
		channelSize:  channelSize,
		channels:     make(map[string]chan string),
		handlersPool: make([]func(context.Context) error, 0),
		mu:           sync.RWMutex{},
	}
}

func (c *Conveyer) makeChannel(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.channels[name]; !ok {
		c.channels[name] = make(chan string, c.channelSize)
	}
}

func (c *Conveyer) addToPool(function func(context.Context) error) {
	c.handlersPool = append(c.handlersPool, function)
}

func (c *Conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, channel := range c.channels {
		close(channel)
	}
}

func (c *Conveyer) RegisterDecorator(
	decoratorFunc func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.makeChannel(input)
	c.makeChannel(output)

	c.addToPool(func(context context.Context) error {
		c.mu.RLock()
		defer c.mu.RUnlock()

		inChan := c.channels[input]
		outChan := c.channels[output]

		return decoratorFunc(context, inChan, outChan)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	multiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	for _, input := range inputs {
		c.makeChannel(input)
	}

	c.makeChannel(output)

	c.addToPool(func(context context.Context) error {
		c.mu.RLock()
		defer c.mu.RUnlock()

		inputChannels := make([]chan string, len(inputs))
		outChan := c.channels[output]

		for i, inputName := range inputs {
			inputChannels[i] = c.channels[inputName]
		}

		return multiplexerFunc(
			context,
			inputChannels,
			outChan,
		)
	})
}

func (c *Conveyer) RegisterSeparator(
	separatorFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.makeChannel(input)

	for _, output := range outputs {
		c.makeChannel(output)
	}

	c.addToPool(func(context context.Context) error {
		c.mu.RLock()
		defer c.mu.RUnlock()

		inChan := c.channels[input]
		outputChannels := make([]chan string, len(outputs))

		for i, outputName := range outputs {
			outputChannels[i] = c.channels[outputName]
		}

		return separatorFunc(
			context,
			inChan,
			outputChannels,
		)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	if len(c.handlersPool) == 0 {
		return errNoHandlers
	}

	handlerGroup, hCtx := errgroup.WithContext(ctx)

	for _, handler := range c.handlersPool {
		handlerGroup.Go(func() error {
			return handler(hCtx)
		})
	}

	if err := handlerGroup.Wait(); err != nil {
		c.closeAllChannels()

		return fmt.Errorf("conveyer handlers failed: %w", err)
	}

	c.closeAllChannels()

	return nil
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.RLock()
	channel, exists := c.channels[input]
	c.mu.RUnlock()

	if !exists {
		return errChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	channel, exists := c.channels[output]
	c.mu.RUnlock()

	if !exists {
		return "", errChanNotFound
	}

	data, ok := <-channel
	if !ok {
		return "undefined", nil
	}

	return data, nil
}
