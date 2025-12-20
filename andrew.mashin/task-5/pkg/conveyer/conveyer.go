package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var errChannel = errors.New("chan not found")

const undefined = "undefined"

type Conveyer struct {
	channelSize int
	channels    map[string]chan string
	handlers    []func(ctx context.Context) error
	mu          sync.RWMutex
}

func (conv *Conveyer) getOrCreateChannel(name string) chan string {
	conv.mu.RLock()
	channel, exists := conv.channels[name]
	conv.mu.RUnlock()

	if exists {
		return channel
	}

	conv.mu.Lock()
	defer conv.mu.Unlock()

	if channel, exists := conv.channels[name]; exists {
		return channel
	}

	channel = make(chan string, conv.channelSize)
	conv.channels[name] = channel

	return channel
}

func New(size int) *Conveyer {
	if size <= 0 {
		size = 0
	}

	return &Conveyer{
		channelSize: size,
		channels:    make(map[string]chan string),
		handlers:    make([]func(ctx context.Context) error, 0),
		mu:          sync.RWMutex{},
	}
}

func (conv *Conveyer) RegisterDecorator(
	callback func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	inputChannel := conv.getOrCreateChannel(input)
	outputChannel := conv.getOrCreateChannel(output)

	conv.mu.Lock()
	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return callback(ctx, inputChannel, outputChannel)
	})
	conv.mu.Unlock()
}

func (conv *Conveyer) RegisterMultiplexer(
	callback func(
		ctx context.Context,
		inputs []chan string,
		output chan string,
	) error,
	inputs []string,
	output string,
) {
	inputChannels := make([]chan string, len(inputs))

	for i, inputName := range inputs {
		inputChannels[i] = conv.getOrCreateChannel(inputName)
	}

	outputCh := conv.getOrCreateChannel(output)

	conv.mu.Lock()
	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return callback(ctx, inputChannels, outputCh)
	})
	conv.mu.Unlock()
}

func (conv *Conveyer) RegisterSeparator(
	callback func(
		ctx context.Context,
		input chan string,
		outputs []chan string,
	) error,
	input string,
	outputs []string,
) {
	inputChannel := conv.getOrCreateChannel(input)
	outputChannels := make([]chan string, len(outputs))

	for i, outputName := range outputs {
		outputChannels[i] = conv.getOrCreateChannel(outputName)
	}

	conv.mu.Lock()
	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return callback(ctx, inputChannel, outputChannels)
	})
	conv.mu.Unlock()
}

func (conv *Conveyer) Run(ctx context.Context) error {
	group, cont := errgroup.WithContext(ctx)

	conv.mu.RLock()
	handlers := make([]func(ctx context.Context) error, len(conv.handlers))
	copy(handlers, conv.handlers)
	conv.mu.RUnlock()

	for _, handler := range handlers {
		group.Go(func() error {
			return handler(cont)
		})
	}

	err := group.Wait()
	if err != nil {
		return fmt.Errorf("conveyer run: %w", err)
	}

	return nil
}

func (conv *Conveyer) Send(input string, data string) error {
	conv.mu.RLock()
	channel, ok := conv.channels[input]
	conv.mu.RUnlock()

	if !ok {
		return errChannel
	}

	channel <- data

	return nil
}

func (conv *Conveyer) Recv(output string) (string, error) {
	conv.mu.RLock()
	channel, ok1 := conv.channels[output]
	conv.mu.RUnlock()

	if !ok1 {
		return "", errChannel
	}

	data, ok2 := <-channel
	if !ok2 {
		return undefined, nil
	}

	return data, nil
}
