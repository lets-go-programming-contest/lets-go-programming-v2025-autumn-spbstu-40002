package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	channels map[string]chan string
	size     int
	handlers []func(ctx context.Context) error
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		size:     size,
		handlers: make([]func(ctx context.Context) error, 0),
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	if channel, exists := c.channels[name]; exists {
		return channel
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel

	return channel
}

func (c *Conveyer) getChannel(name string) (chan string, error) {
	channel, exists := c.channels[name]
	if !exists {
		return nil, ErrChanNotFound
	}

	return channel, nil
}

func (c *Conveyer) RegisterDecorator(
	decoratorFunc func(ctx context.Context, input chan string, output chan string) error,
	inputName string,
	outputName string,
) {
	c.getOrCreateChannel(inputName)
	c.getOrCreateChannel(outputName)

	handler := func(ctx context.Context) error {
		input := c.getOrCreateChannel(inputName)
		output := c.getOrCreateChannel(outputName)

		return decoratorFunc(ctx, input, output)
	}

	c.handlers = append(c.handlers, handler)
}

func (c *Conveyer) RegisterMultiplexer(
	multiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputNames []string,
	outputName string,
) {
	for _, name := range inputNames {
		c.getOrCreateChannel(name)
	}

	c.getOrCreateChannel(outputName)

	handler := func(ctx context.Context) error {
		inputs := make([]chan string, len(inputNames))
		for i, name := range inputNames {
			inputs[i] = c.getOrCreateChannel(name)
		}

		output := c.getOrCreateChannel(outputName)

		return multiplexerFunc(ctx, inputs, output)
	}

	c.handlers = append(c.handlers, handler)
}

func (c *Conveyer) RegisterSeparator(
	separatorFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	inputName string,
	outputNames []string,
) {
	c.getOrCreateChannel(inputName)

	for _, name := range outputNames {
		c.getOrCreateChannel(name)
	}

	handler := func(ctx context.Context) error {
		input := c.getOrCreateChannel(inputName)
		outputs := make([]chan string, len(outputNames))

		for i, name := range outputNames {
			outputs[i] = c.getOrCreateChannel(name)
		}

		return separatorFunc(ctx, input, outputs)
	}

	c.handlers = append(c.handlers, handler)
}

func (c *Conveyer) Send(inputName string, data string) error {
	channel, err := c.getChannel(inputName)
	if err != nil {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(outputName string) (string, error) {
	channel, err := c.getChannel(outputName)
	if err != nil {
		return "", ErrChanNotFound
	}

	data, ok := <-channel
	if !ok {
		return "undefined", nil
	}

	return data, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	waitGroup := sync.WaitGroup{}
	errorChannel := make(chan error, len(c.handlers))

	for _, currentHandler := range c.handlers {
		waitGroup.Add(1)

		handlerCopy := currentHandler
		go func() {
			defer waitGroup.Done()

			if err := handlerCopy(ctx); err != nil {
				select {
				case errorChannel <- err:
				default:
				}
			}
		}()
	}

	go func() {
		waitGroup.Wait()
		close(errorChannel)
	}()

	select {
	case err := <-errorChannel:
		return err
	case <-ctx.Done():
		waitGroup.Wait()

		return nil
	}
}