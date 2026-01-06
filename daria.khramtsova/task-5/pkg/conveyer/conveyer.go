package conveyer

import (
	"context"
	"errors"
	"fmt"
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

func (conveyer *Conveyer) getOrCreate(name string) chan string {
	conveyer.mu.RLock()
	channel, exists := conveyer.channels[name]
	conveyer.mu.RUnlock()

	if exists {
		return channel
	}

	conveyer.mu.Lock()
	defer conveyer.mu.Unlock()

	channel, exists = conveyer.channels[name]
	if exists {
		return channel
	}

	channel = make(chan string, conveyer.channelSize)
	conveyer.channels[name] = channel

	return channel
}

func (conveyer *Conveyer) closeChannels() {
	for _, channel := range conveyer.channels {
		close(channel)
	}
}

func (conveyer *Conveyer) RegisterDecorator(
	handlerFunc func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	inputChan := conveyer.getOrCreate(input)
	outputChan := conveyer.getOrCreate(output)

	conveyer.mu.Lock()
	conveyer.handlers = append(conveyer.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChan, outputChan)
	})
	conveyer.mu.Unlock()
}

func (conveyer *Conveyer) RegisterMultiplexer(
	handlerFunc func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	inputChannels := make([]chan string, len(inputs))
	for index, name := range inputs {
		inputChannels[index] = conveyer.getOrCreate(name)
	}

	outputChan := conveyer.getOrCreate(output)

	conveyer.mu.Lock()
	conveyer.handlers = append(conveyer.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannels, outputChan)
	})
	conveyer.mu.Unlock()
}

func (conveyer *Conveyer) RegisterSeparator(
	handlerFunc func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	inputChan := conveyer.getOrCreate(input)

	outputChannels := make([]chan string, len(outputs))
	for index, name := range outputs {
		outputChannels[index] = conveyer.getOrCreate(name)
	}

	conveyer.mu.Lock()
	conveyer.handlers = append(conveyer.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChan, outputChannels)
	})
	conveyer.mu.Unlock()
}

func (conveyer *Conveyer) Run(ctx context.Context) error {
	group, groupCtx := errgroup.WithContext(ctx)

	conveyer.mu.RLock()
	handlersCopy := append([]func(context.Context) error{}, conveyer.handlers...)
	conveyer.mu.RUnlock()

	for _, handler := range handlersCopy {
		handler := handler
		group.Go(func() error {
			return handler(groupCtx)
		})
	}

	err := group.Wait()

	conveyer.mu.Lock()
	conveyer.closeChannels()
	conveyer.mu.Unlock()

	if err != nil {
		return fmt.Errorf("conveyer run: %w", err)
	}

	return nil
}

func (conveyer *Conveyer) Send(name string, data string) error {
	conveyer.mu.RLock()
	channel, exists := conveyer.channels[name]
	conveyer.mu.RUnlock()

	if !exists {
		return errChanNotFound
	}

	channel <- data

	return nil
}

func (conveyer *Conveyer) Recv(name string) (string, error) {
	conveyer.mu.RLock()
	channel, exists := conveyer.channels[name]
	conveyer.mu.RUnlock()

	if !exists {
		return "", errChanNotFound
	}

	data, isOpen := <-channel
	if !isOpen {
		return "undefined", nil
	}

	return data, nil
}
