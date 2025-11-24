package conveyer

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

const undefinedData = "undefined"

var ErrChanNotFound = errors.New("chan not found")

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

func (c *Conveyer) makeChannel(name string) {
	if _, ok := c.channels[name]; !ok {
		c.channels[name] = make(chan string, c.channelSize)
	}
}

func (c *Conveyer) makeChannels(names ...string) {
	for _, n := range names {
		c.makeChannel(n)
	}
}

func (c *Conveyer) RegisterDecorator(
	handler func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.makeChannels(input)
	c.makeChannels(output)

	c.handlersPool = append(c.handlersPool, func(ctx context.Context) error {

		return handler(ctx, c.channels[input], c.channels[output])
	})
}

func (c *Conveyer) RegisterMultiplexer(
	handler func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.makeChannels(inputs...)
	c.makeChannels(output)

	c.handlersPool = append(c.handlersPool, func(ctx context.Context) error {
		ins := make([]chan string, 0, len(inputs))
		for _, n := range inputs {
			ins = append(ins, c.channels[n])
		}

		return handler(ctx, ins, c.channels[output])
	})
}

func (c *Conveyer) RegisterSeparator(
	handler func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.makeChannels(input)
	c.makeChannels(outputs...)

	c.handlersPool = append(c.handlersPool, func(ctx context.Context) error {
		outs := make([]chan string, 0, len(outputs))
		for _, n := range outputs {
			outs = append(outs, c.channels[n])
		}

		return handler(ctx, c.channels[input], outs)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	errGroup, egCtx := errgroup.WithContext(ctx)

	for _, handlerFunc := range c.handlersPool {
		hf := handlerFunc

		errGroup.Go(func() error {
			return hf(egCtx)
		})
	}

	return errGroup.Wait()
}

func (c *Conveyer) Send(input string, data string) error {
	ch, ok := c.channels[input]
	if !ok {
		return ErrChanNotFound
	}

	ch <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	ch, ok := c.channels[output]
	if !ok {
		return "", ErrChanNotFound
	}

	v, open := <-ch

	if !open {
		return undefinedData, nil
	}

	return v, nil
}
