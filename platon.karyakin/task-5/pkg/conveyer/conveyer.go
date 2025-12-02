package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

const undefined = "undefined"

type conveyerImpl struct {
	size         int
	channels     map[string]chan string
	mu           sync.RWMutex
	decorators   []decoratorEntry
	multiplexers []multiplexerEntry
	separators   []separatorEntry
}

type decoratorEntry struct {
	function func(ctx context.Context, input chan string, output chan string) error
	input    string
	output   string
}

type multiplexerEntry struct {
	function func(ctx context.Context, inputs []chan string, output chan string) error
	inputs   []string
	output   string
}

type separatorEntry struct {
	function func(ctx context.Context, input chan string, outputs []chan string) error
	input    string
	outputs  []string
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		size:         size,
		channels:     make(map[string]chan string),
		mu:           sync.RWMutex{},
		decorators:   make([]decoratorEntry, 0),
		multiplexers: make([]multiplexerEntry, 0),
		separators:   make([]separatorEntry, 0),
	}
}

func (c *conveyerImpl) getOrCreateChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch

	return ch
}

func (c *conveyerImpl) getChan(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, exists := c.channels[name]

	return ch, exists
}

func (c *conveyerImpl) Send(input string, data string) error {
	ch, exists := c.getChan(input)
	if !exists {
		return ErrChanNotFound
	}

	ch <- data

	return nil
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	ch, exists := c.getChan(output)
	if !exists {
		return "", ErrChanNotFound
	}

	data, ok := <-ch
	if !ok {
		return undefined, nil
	}

	return data, nil
}

func (c *conveyerImpl) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *conveyerImpl) RegisterDecorator(
	function func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.getOrCreateChan(input)
	c.getOrCreateChan(output)

	c.decorators = append(c.decorators, decoratorEntry{
		function: function,
		input:    input,
		output:   output,
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	function func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	for _, name := range inputs {
		c.getOrCreateChan(name)
	}

	c.getOrCreateChan(output)

	c.multiplexers = append(c.multiplexers, multiplexerEntry{
		function: function,
		inputs:   inputs,
		output:   output,
	})
}

func (c *conveyerImpl) RegisterSeparator(
	function func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.getOrCreateChan(input)

	for _, name := range outputs {
		c.getOrCreateChan(name)
	}

	c.separators = append(c.separators, separatorEntry{
		function: function,
		input:    input,
		outputs:  outputs,
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	defer c.closeAllChannels()

	group, groupCtx := errgroup.WithContext(ctx)

	for _, decorator := range c.decorators {
		inputChan := c.getOrCreateChan(decorator.input)
		outputChan := c.getOrCreateChan(decorator.output)

		group.Go(func() error {
			return decorator.function(groupCtx, inputChan, outputChan)
		})
	}

	for _, multiplexer := range c.multiplexers {
		outputChan := c.getOrCreateChan(multiplexer.output)
		inputChannels := make([]chan string, len(multiplexer.inputs))

		for index, name := range multiplexer.inputs {
			inputChannels[index] = c.getOrCreateChan(name)
		}

		group.Go(func() error {
			return multiplexer.function(groupCtx, inputChannels, outputChan)
		})
	}

	for _, separator := range c.separators {
		inputChan := c.getOrCreateChan(separator.input)
		outputChannels := make([]chan string, len(separator.outputs))

		for index, name := range separator.outputs {
			outputChannels[index] = c.getOrCreateChan(name)
		}

		group.Go(func() error {
			return separator.function(groupCtx, inputChan, outputChannels)
		})
	}

	err := group.Wait()
	if err != nil {
		return fmt.Errorf("conveyer run failed: %w", err)
	}

	return nil
}
