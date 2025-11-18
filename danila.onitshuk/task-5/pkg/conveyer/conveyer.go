package conveyer

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Conveyer struct {
	channelSize  int
	channels     map[string]chan string
	handlersPool []func(context.Context) error
}

func New(channelSize int) *Conveyer {
	return &Conveyer{
		channelSize:  channelSize,
		channels:     make(map[string]chan string),
		handlersPool: make([]func(context.Context) error, 0),
	}
}

// Добавление декоратора.
func (c *Conveyer) RegisterDecorator(
	fn func(
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
		return fn(ctx, c.channels[input], c.channels[output])
	})
}

// Добавление мультиплексора.
func (c *Conveyer) RegisterMultiplexer(
	fn func(
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

		return fn(ctx, inputsCh, c.channels[output])
	})
}

// Добавление сепаратора.
func (c *Conveyer) RegisterSeparator(
	fn func(
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

		return fn(ctx, c.channels[input], outputsCh)
	})
}

// Запуск конвеера.
func (c *Conveyer) Run(ctx context.Context) error {
	defer c.closeChannels()

	groupHendlers, gctx := errgroup.WithContext(ctx)

	for _, hendler := range c.handlersPool {
		h := hendler
		groupHendlers.Go(func() error {
			return h(gctx)
		})
	}

	return groupHendlers.Wait()
}

// Отправка сообщение в канал input.
func (c *Conveyer) Send(input string, data string) error {
	ch, ok := c.channels[input]
	if !ok {
		return ErrNoChannel
	}

	ch <- data

	return nil
}

// Получение данных из канала output.
func (c *Conveyer) Recv(output string) (string, error) {
	ch, ok := c.channels[output]

	if !ok {
		return "", ErrNoChannel
	}

	data, ok := <-ch
	if !ok {
		return undefinedData, nil
	}

	return data, nil
}
