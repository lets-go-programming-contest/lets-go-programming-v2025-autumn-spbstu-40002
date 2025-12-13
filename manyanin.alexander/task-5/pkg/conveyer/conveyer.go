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
	conv.mu.Lock()

	defer conv.mu.Unlock()

	if ch, exists := conv.channels[name]; exists {
		return ch
	}

	ch := make(chan string, conv.channelSize)

	conv.channels[name] = ch

	return ch
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
	decoratorFunc func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	inputChannel := conv.getOrCreateChannel(input)
	outputChannel := conv.getOrCreateChannel(output)
	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return decoratorFunc(ctx, inputChannel, outputChannel)
	})
}

func (conv *Conveyer) RegisterMultiplexer(
	multiplexerFunc func(
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
	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return multiplexerFunc(ctx, inputChannels, outputCh)
	})
}

func (conv *Conveyer) RegisterSeparator(
	separatorFunc func(
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

	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return separatorFunc(ctx, inputChannel, outputChannels)
	})
}

func (conv *Conveyer) Run(ctx context.Context) error {
	group, cont := errgroup.WithContext(ctx)

	for _, handler := range conv.handlers {
		currentHandler := handler

		group.Go(func() error {
			return currentHandler(cont)
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

	channel, normal := conv.channels[output]

	conv.mu.RUnlock()

	if !normal {
		return "", errChannel
	}

	data, normal := <-channel
	if !normal {
		return undefined, nil
	}

	return data, nil
}

func (conv *Conveyer) Close() {
	conv.mu.Lock()

	defer conv.mu.Unlock()

	for name, ch := range conv.channels {
		select {
		case <-ch:
		default:
		}
		close(ch)
		delete(conv.channels, name)
	}
}
