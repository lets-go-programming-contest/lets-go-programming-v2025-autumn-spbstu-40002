package conveyer

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

const (
	undefinedData = "undefined"
)

type Conveyer struct {
	channelSize  int
	channels     map[string]chan string
	handlersPool []func(context.Context) error
	mutex        sync.RWMutex
}

func New(channelSize int) *Conveyer {
	return &Conveyer{
		channelSize:  channelSize,
		channels:     make(map[string]chan string),
		handlersPool: make([]func(context.Context) error, 0),
		mutex:        sync.RWMutex{},
	}
}

// Добавление декоратора.
func (c *Conveyer) RegisterDecorator(
	hendlerFunc func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	c.makeChannels(input)
	c.makeChannels(output)
	c.handlersPool = append(c.handlersPool, func(ctx context.Context) error {
		return hendlerFunc(ctx, c.channels[input], c.channels[output])
	})
}

// Добавление мультиплексора.
func (c *Conveyer) RegisterMultiplexer(
	hendlerFunc func(
		ctx context.Context,
		inputs []chan string,
		output chan string,
	) error,
	inputs []string,
	output string,
) {
	c.makeChannels(inputs...)
	c.makeChannels(output)
	c.handlersPool = append(c.handlersPool, func(ctx context.Context) error {
		inputsCh := make([]chan string, 0, len(inputs))

		for _, ch := range inputs {
			inputsCh = append(inputsCh, c.channels[ch])
		}

		return hendlerFunc(ctx, inputsCh, c.channels[output])
	})
}

// Добавление сепаратора.
func (c *Conveyer) RegisterSeparator(
	hendlerFunc func(
		ctx context.Context,
		input chan string,
		outputs []chan string,
	) error,
	input string,
	outputs []string,
) {
	c.makeChannels(input)
	c.makeChannels(outputs...)
	c.handlersPool = append(c.handlersPool, func(ctx context.Context) error {
		outputsCh := make([]chan string, 0, len(outputs))

		for _, ch := range outputs {
			outputsCh = append(outputsCh, c.channels[ch])
		}

		return hendlerFunc(ctx, c.channels[input], outputsCh)
	})
}

// Запуск конвеера.
func (c *Conveyer) Run(ctx context.Context) error {
	groupHendlers, gctx := errgroup.WithContext(ctx)

	for _, hendler := range c.handlersPool {
		h := hendler

		groupHendlers.Go(func() error {
			return h(gctx)
		})
	}

	return groupHendlers.Wait() //nolint:wrapcheck
}

// Отправка сообщение в канал input.
func (c *Conveyer) Send(input string, data string) error {
	c.mutex.RLock()
	channel, ok := c.channels[input]
	c.mutex.RUnlock()

	if !ok {
		return ErrNoChannel
	}

	channel <- data

	return nil
}

// Получение данных из канала output.
func (c *Conveyer) Recv(output string) (string, error) {
	c.mutex.RLock()
	channel, channelAvailability := c.channels[output]
	c.mutex.RUnlock()

	if !channelAvailability {
		return "", ErrNoChannel
	}

	data, dataAvailability := <-channel
	if !dataAvailability {
		return undefinedData, nil
	}

	return data, nil
}
