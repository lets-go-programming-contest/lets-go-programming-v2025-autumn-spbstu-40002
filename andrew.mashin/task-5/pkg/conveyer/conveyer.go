package conveyer

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

var errChannel = errors.New("chan not found")

const undefined = "undefined"

type Conveyer struct {
	channelSize int
	channels    map[string]chan string
	handlers    []func(ctx context.Context) error
}

func (conv *Conveyer) getOrCreateChannel(name string) chan string {
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
	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return callback(ctx, inputChannel, outputChannel)
	})
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
	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return callback(ctx, inputChannels, outputCh)
	})
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

	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return callback(ctx, inputChannel, outputChannels)
	})
}

func (conv *Conveyer) Run(ctx context.Context) error {
	group, cont := errgroup.WithContext(ctx)

	for _, handler := range conv.handlers {
		group.Go(func() error {
			return handler(cont)
		})
	}

	return fmt.Errorf("group wait: %w", group.Wait())
}

func (conv *Conveyer) Send(input string, data string) error {
	channel, ok := conv.channels[input]
	if !ok {
		return errChannel
	}

	channel <- data

	return nil
}

func (conv *Conveyer) Recv(output string) (string, error) {
	channel, ok1 := conv.channels[output]
	if !ok1 {
		return "", errChannel
	}

	data, ok2 := <-channel
	if !ok2 {
		return undefined, nil
	}

	return data, nil
}
