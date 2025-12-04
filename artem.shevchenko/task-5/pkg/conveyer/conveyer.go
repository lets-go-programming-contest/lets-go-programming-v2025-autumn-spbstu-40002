package conveyer

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

const undefinedData = "undefined"

var ErrChanNotFound = errors.New("chan not found")

type impl struct {
	bufferSize int
	channels   map[string]chan string
	tasks      []func(ctx context.Context) error
}

func New(bufferSize int) conveyer {
	return &impl{
		bufferSize: bufferSize,
		channels:   make(map[string]chan string),
		tasks:      make([]func(context.Context) error, 0),
	}
}

func (conveyor *impl) getOrCreateChan(name string) chan string {
	if channel, exists := conveyor.channels[name]; exists {
		return channel
	}
	channel := make(chan string, conveyor.bufferSize)
	conveyor.channels[name] = channel
	return channel
}

func (conveyor *impl) RegisterDecorator(
	handler func(ctx context.Context, inputChannel, outputChannel chan string) error,
	inputName, outputName string,
) {
	inputChannel := conveyor.getOrCreateChan(inputName)
	outputChannel := conveyor.getOrCreateChan(outputName)
	conveyor.tasks = append(conveyor.tasks, func(ctx context.Context) error {
		return handler(ctx, inputChannel, outputChannel)
	})
}

func (conveyor *impl) RegisterMultiplexer(
	handler func(ctx context.Context, inputChannels []chan string, outputChannel chan string) error,
	inputNames []string,
	outputName string,
) {
	inputChannels := make([]chan string, len(inputNames))
	for index, name := range inputNames {
		inputChannels[index] = conveyor.getOrCreateChan(name)
	}
	outputChannel := conveyor.getOrCreateChan(outputName)
	conveyor.tasks = append(conveyor.tasks, func(ctx context.Context) error {
		return handler(ctx, inputChannels, outputChannel)
	})
}

func (conveyor *impl) RegisterSeparator(
	handler func(ctx context.Context, inputChannel chan string, outputChannels []chan string) error,
	inputName string,
	outputNames []string,
) {
	inputChannel := conveyor.getOrCreateChan(inputName)
	outputChannels := make([]chan string, len(outputNames))
	for index, name := range outputNames {
		outputChannels[index] = conveyor.getOrCreateChan(name)
	}
	conveyor.tasks = append(conveyor.tasks, func(ctx context.Context) error {
		return handler(ctx, inputChannel, outputChannels)
	})
}

func (conveyor *impl) Run(ctx context.Context) error {
	if len(conveyor.tasks) == 0 {
		<-ctx.Done()
		return ctx.Err()
	}

	group, ctx := errgroup.WithContext(ctx)

	for _, task := range conveyor.tasks {
		task := task
		group.Go(func() error {
			return task(ctx)
		})
	}

	if err := group.Wait(); err != nil {
		for _, channel := range conveyor.channels {
			close(channel)
		}
		return err
	}

	for _, channel := range conveyor.channels {
		close(channel)
	}
	return nil
}

func (conveyor *impl) Send(channelName, data string) error {
	channel, exists := conveyor.channels[channelName]
	if !exists {
		return ErrChanNotFound
	}
	channel <- data
	return nil
}

func (conveyor *impl) Recv(channelName string) (string, error) {
	channel, exists := conveyor.channels[channelName]
	if !exists {
		return "", ErrChanNotFound
	}
	data, channelOpen := <-channel
	if !channelOpen {
		return undefinedData, nil
	}
	return data, nil
}