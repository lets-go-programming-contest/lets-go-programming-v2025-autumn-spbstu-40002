package conveyer

import (
	"context"
	"errors"
)

var undefinedChannel = errors.New("undefined")
var nonExistingСhannel = errors.New("chan not found")
var contextIsCanceled = errors.New("the context is canceled")

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
	fn func(
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
		return fn(ctx, conveyer.mapChannels[input], conveyer.mapChannels[output])
	})
}

func (conveyer *Conveyer) RegisterMultiplexer(
	fn func(
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

		return fn(ctx, inputsCh, conveyer.mapChannels[output])
	})
}

func (conveyer *Conveyer) RegisterSeparator(
	fn func(
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
		return fn(ctx, conveyer.mapChannels[input], outputsCh)
	})
}

func (conveyer *Conveyer) Run(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		for _, channel := range conveyer.mapChannels {
			close(channel)
		}
		return contextIsCanceled
	}
	for _, fn := range conveyer.handlersPool {
		go fn(ctx)
	}
	return nil
}

func (conveyer *Conveyer) Send(input string, data string) error {
	if _, ok := conveyer.mapChannels[input]; !ok {
		return nonExistingСhannel
	} else {
		conveyer.mapChannels[input] <- data
		return nil
	}
}

func (conveyer *Conveyer) Recv(output string) (string, error) {
	channel, ok := conveyer.mapChannels[output]
	if !ok {
		return "", nonExistingСhannel
	} else {
		var outputString string
		outputString, ok = <-channel
		if ok {
			return outputString, nil
		} else {
			return "", undefinedChannel
		}
	}
}
