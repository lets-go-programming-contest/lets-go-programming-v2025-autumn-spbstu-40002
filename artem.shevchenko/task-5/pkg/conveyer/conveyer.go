package conveyer

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Conveyer struct {
	mutex        sync.RWMutex
	channelSize  int
	channels     map[string]chan string
	handlersPool []func(context.Context) error
}

func New(channelSize int) *Conveyer {
	return &Conveyer{
		mutex:        sync.RWMutex{},
		channelSize:  channelSize,
		channels:     make(map[string]chan string),
		handlersPool: make([]func(context.Context) error, 0),
	}
}

func (conveyer *Conveyer) makeChannels(names ...string) {
	conveyer.mutex.Lock()
	defer conveyer.mutex.Unlock()

	for _, name := range names {
		if _, channelExists := conveyer.channels[name]; !channelExists {
			conveyer.channels[name] = make(chan string, conveyer.channelSize)
		}
	}
}

func (conveyer *Conveyer) RegisterDecorator(
	task func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	conveyer.makeChannels(input, output)

	conveyer.mutex.Lock()
	conveyer.handlersPool = append(conveyer.handlersPool, func(ctx context.Context) error {
		return task(ctx, conveyer.channels[input], conveyer.channels[output])
	})
	conveyer.mutex.Unlock()
}

func (conveyer *Conveyer) RegisterMultiplexer(
	task func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	conveyer.makeChannels(inputs...)
	conveyer.makeChannels(output)

	conveyer.mutex.Lock()
	conveyer.handlersPool = append(conveyer.handlersPool, func(ctx context.Context) error {
		inputChannels := make([]chan string, len(inputs))
		for index, input := range inputs {
			inputChannels[index] = conveyer.channels[input]
		}

		return task(ctx, inputChannels, conveyer.channels[output])
	})
	conveyer.mutex.Unlock()
}

func (conveyer *Conveyer) RegisterSeparator(
	task func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	conveyer.makeChannels(input)
	conveyer.makeChannels(outputs...)

	conveyer.mutex.Lock()
	conveyer.handlersPool = append(conveyer.handlersPool, func(ctx context.Context) error {
		outputChannels := make([]chan string, len(outputs))
		for index, output := range outputs {
			outputChannels[index] = conveyer.channels[output]
		}

		return task(ctx, conveyer.channels[input], outputChannels)
	})
	conveyer.mutex.Unlock()
}

func (conveyer *Conveyer) Run(ctx context.Context) error {
	handlersGroup, handlersContext := errgroup.WithContext(ctx)

	for _, handler := range conveyer.handlersPool {
		currentHandler := handler

		handlersGroup.Go(func() error {
			return currentHandler(handlersContext)
		})
	}

	if err := handlersGroup.Wait(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (conveyer *Conveyer) Send(input string, data string) error {
	conveyer.mutex.RLock()
	channel, channelExists := conveyer.channels[input]
	conveyer.mutex.RUnlock()

	if !channelExists {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (conveyer *Conveyer) Recv(output string) (string, error) {
	conveyer.mutex.RLock()
	channel, channelExists := conveyer.channels[output]
	conveyer.mutex.RUnlock()

	if !channelExists {
		return "", ErrChanNotFound
	}

	data, channelOpen := <-channel
	if !channelOpen {
		return UndefinedData, nil
	}

	return data, nil
}
