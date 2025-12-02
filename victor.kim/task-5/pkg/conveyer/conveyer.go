package conveyer

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

const undefinedData = "undefined"

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	channelSize int
	channels    map[string]chan string
	handlers    []func(context.Context) error
}

func New(channelSize int) *Conveyer {
	return &Conveyer{
		channelSize: channelSize,
		channels:    make(map[string]chan string),
		handlers:    make([]func(context.Context) error, 0),
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
	out string,
) {
	c.makeChannels(input)
	c.makeChannels(out)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, c.channels[input], c.channels[out])
	})
}

func (c *Conveyer) RegisterMultiplexer(
	handler func(ctx context.Context, inputs []chan string, output chan string) error,
	inNames []string,
	out string,
) {
	c.makeChannels(inNames...)
	c.makeChannels(out)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		inputChans := make([]chan string, 0, len(inNames))
		for _, name := range inNames {
			inputChans = append(inputChans, c.channels[name])
		}

		return handler(ctx, inputChans, c.channels[out])
	})
}

func (c *Conveyer) RegisterSeparator(
	handler func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outNames []string,
) {
	c.makeChannels(input)
	c.makeChannels(outNames...)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		outputChans := make([]chan string, 0, len(outNames))
		for _, name := range outNames {
			outputChans = append(outputChans, c.channels[name])
		}

		return handler(ctx, c.channels[input], outputChans)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	errGroup, egCtx := errgroup.WithContext(ctx)

	for _, handlerFunc := range c.handlers {
		hf := handlerFunc

		errGroup.Go(func() error {
			return hf(egCtx)
		})
	}

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("conveyer handlers: %w", err)
	}

	return nil
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

	value, open := <-ch

	if !open {
		return undefinedData, nil
	}

	return value, nil
}
