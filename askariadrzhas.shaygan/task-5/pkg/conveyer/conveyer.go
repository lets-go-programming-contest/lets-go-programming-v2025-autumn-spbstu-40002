package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var errChannelNotFound = errors.New("channel does not exist")

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
	channels     map[string]chan string
	bufferLength int
	processors   []func(ctx context.Context) error
	lock         sync.RWMutex
}

func New(size int) Conveyer {
	return Conveyer{
		channels:     make(map[string]chan string),
		bufferLength: size,
		processors:   make([]func(ctx context.Context) error, 0),
		lock:         sync.RWMutex{},
	}
}

func (c *Conveyer) setupChannel(channelName string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, exists := c.channels[channelName]; !exists {
		c.channels[channelName] = make(chan string, c.bufferLength)
	}
}

func (c *Conveyer) fetchChannel(channelName string) (chan string, error) {
	if ch, ok := c.channels[channelName]; ok {
		return ch, nil
	}

	return nil, errChannelNotFound
}

func (c *Conveyer) terminateChannels() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *Conveyer) addProcessor(processor func(ctx context.Context) error) {
	c.processors = append(c.processors, processor)
}

func (c *Conveyer) RegisterDecorator(
	function func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string, output string,
) {
	c.setupChannel(input)
	c.setupChannel(output)

	c.addProcessor(func(ctx context.Context) error {
		c.lock.RLock()
		defer c.lock.RUnlock()

		inputCh, _ := c.fetchChannel(input)
		outputCh, _ := c.fetchChannel(output)

		return function(ctx, inputCh, outputCh)
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
		c.setupChannel(inp)
	}

	c.setupChannel(output)

	c.addProcessor(func(ctx context.Context) error {
		c.lock.RLock()
		defer c.lock.RUnlock()

		inputChannels := make([]chan string, len(inputs))

		for i, inp := range inputs {
			ch, _ := c.fetchChannel(inp)
			inputChannels[i] = ch
		}

		outputCh, _ := c.fetchChannel(output)

		return function(ctx, inputChannels, outputCh)
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
	c.setupChannel(input)

	for _, out := range outputs {
		c.setupChannel(out)
	}

	c.addProcessor(func(ctx context.Context) error {
		c.lock.RLock()
		defer c.lock.RUnlock()

		inputCh, _ := c.fetchChannel(input)

		outputChannels := make([]chan string, len(outputs))

		for i, out := range outputs {
			ch, _ := c.fetchChannel(out)
			outputChannels[i] = ch
		}

		return function(ctx, inputCh, outputChannels)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	defer c.terminateChannels()

	group, ctxWithCancel := errgroup.WithContext(ctx)

	c.lock.RLock()

	for _, proc := range c.processors {
		currentProc := proc
		group.Go(func() error {
			return currentProc(ctxWithCancel)
		})
	}

	c.lock.RUnlock()

	if err := group.Wait(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(input string, data string) error {
	c.lock.RLock()

	inputCh, err := c.fetchChannel(input)
	if err != nil {
		return err
	}

	c.lock.RUnlock()

	inputCh <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.lock.RLock()

	outputCh, err := c.fetchChannel(output)
	if err != nil {
		return "", err
	}

	c.lock.RUnlock()

	if data, ok := <-outputCh; ok {
		return data, nil
	} else {
		return undefined, nil
	}
}
