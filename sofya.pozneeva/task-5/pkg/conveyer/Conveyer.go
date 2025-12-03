package conveyer

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

var errUndefinedChannel = errors.New("channel is undefined")
var errNonExistingChannel = errors.New("chan not found")

func (conveyer *Conveyer) makeChannel(name string) {
	if _, ok := conveyer.mapChannels[name]; !ok {
		conveyer.mapChannels[name] = make(chan string, conveyer.channelSize)
	}
}

func (c *Conveyer) makeChannels(names ...string) {
	for _, name := range names {
		c.makeChannel(name)
	}
}

type Conveyer struct {
	channelSize  int
	mapChannels  map[string]chan string
	handlersPool []func(context.Context) error
}

func (conveyer *Conveyer) RegisterDecorator(
	fnDecorator func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	conveyer.makeChannel(input)
	conveyer.makeChannel(output)
	conveyer.handlersPool = append(conveyer.handlersPool, func(ctx context.Context) error {
		return fnDecorator(ctx, conveyer.mapChannels[input], conveyer.mapChannels[output])
	})
}

func (conveyer *Conveyer) RegisterMultiplexer(
	fnMultiplexer func(
		ctx context.Context,
		inputs []chan string,
		output chan string,
	) error,
	inputs []string,
	output string,
) {
	conveyer.makeChannels(inputs...)
	conveyer.makeChannel(output)
	conveyer.handlersPool = append(conveyer.handlersPool, func(ctx context.Context) error {
		inputsCh := make([]chan string, 0, len(inputs))

		for _, ch := range inputs {
			inputsCh = append(inputsCh, conveyer.mapChannels[ch])
		}

		return fnMultiplexer(ctx, inputsCh, conveyer.mapChannels[output])
	})
}

func (conveyer *Conveyer) RegisterSeparator(
	fnSeparator func(
		ctx context.Context,
		input chan string,
		outputs []chan string,
	) error,
	input string,
	outputs []string,
) {
	conveyer.makeChannel(input)
	conveyer.makeChannels(outputs...)
	conveyer.handlersPool = append(conveyer.handlersPool, func(ctx context.Context) error {
		outputsCh := make([]chan string, 0, len(outputs))
		for _, ch := range outputs {
			outputsCh = append(outputsCh, conveyer.mapChannels[ch])
		}

		return fnSeparator(ctx, conveyer.mapChannels[input], outputsCh)
	})
}

func (conveyer *Conveyer) Run(ctx context.Context) error {
	groupHendlers, gctx := errgroup.WithContext(ctx)

	for _, hendler := range conveyer.handlersPool {
		h := hendler

		groupHendlers.Go(func() error {
			return h(gctx)
		})
	}

	return groupHendlers.Wait() //nolint:wrapcheck
}

func (conveyer *Conveyer) Send(input string, data string) error {
	if _, ok := conveyer.mapChannels[input]; !ok {
		return errNonExistingChannel
	} else {
		conveyer.mapChannels[input] <- data

		return nil
	}
}

func (c *Conveyer) Recv(output string) (string, error) {
	channel, okChan := c.mapChannels[output]
	if !okChan {
		return "", errNonExistingChannel
	} else {
		var outputString string

		outputString, okChan := <-channel
		if okChan {
			return outputString, nil
		} else {
			return "", errUndefinedChannel
		}
	}
}
