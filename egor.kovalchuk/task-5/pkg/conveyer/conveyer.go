package conveyer

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

var (
	errSendChanNotFound = errors.New("conveyer.Send: chan not found")
	errRecvChanNotFound = errors.New("conveyer.Recv: chan not found")
)

type Decorator func(
	context.Context,
	chan string,
	chan string,
) error

type Multiplexer func(
	context.Context,
	[]chan string,
	chan string,
) error

type Separator func(
	context.Context,
	chan string,
	[]chan string,
) error

type Conveyer struct {
	channelSize int
	channels    map[string]chan string
	pool        []func(context.Context) error
}

func New(channelSize int) *Conveyer {
	return &Conveyer{
		channelSize,
		map[string]chan string{},
		[]func(context.Context) error{},
	}
}

func (conveyer *Conveyer) addToPool(function func(context.Context) error) {
	conveyer.pool = append(conveyer.pool, function)
}

func (conveyer *Conveyer) makeChannel(name string) {
	if _, ok := conveyer.channels[name]; !ok {
		conveyer.channels[name] = make(chan string, conveyer.channelSize)
	}
}

func (conveyer *Conveyer) makeChannels(names ...string) {
	for _, name := range names {
		conveyer.makeChannel(name)
	}
}

func (conveyer *Conveyer) collectChannels(names ...string) []chan string {
	result := []chan string{}
	for _, name := range names {
		result = append(result, conveyer.channels[name])
	}

	return result
}

func (conveyer *Conveyer) RegisterDecorator(
	decorator Decorator,
	input string,
	output string,
) {
	conveyer.makeChannels(input, output)
	conveyer.addToPool(func(context context.Context) error {
		return decorator(
			context,
			conveyer.channels[input],
			conveyer.channels[output],
		)
	})
}

func (conveyer *Conveyer) RegisterMultiplexer(
	multiplexer Multiplexer,
	inputs []string,
	output string,
) {
	conveyer.makeChannel(output)
	conveyer.makeChannels(inputs...)
	conveyer.addToPool(func(context context.Context) error {
		return multiplexer(
			context,
			conveyer.collectChannels(inputs...),
			conveyer.channels[output],
		)
	})
}

func (conveyer *Conveyer) RegisterSeparator(
	separator Separator,
	input string,
	outputs []string,
) {
	conveyer.makeChannel(input)
	conveyer.makeChannels(outputs...)
	conveyer.addToPool(func(context context.Context) error {
		return separator(
			context,
			conveyer.channels[input],
			conveyer.collectChannels(outputs...),
		)
	})
}

func (conveyer *Conveyer) Run(context context.Context) error {
	group, contextWithErrs := errgroup.WithContext(context)
	for _, function := range conveyer.pool {
		group.Go(func() error {
			return function(contextWithErrs)
		})
	}

	/* Errors are mine anyway.  */
	return group.Wait() //nolint:wrapcheck
}

func (conveyer *Conveyer) Send(input string, data string) error {
	channel, ok := conveyer.channels[input]
	if !ok {
		return errSendChanNotFound
	}

	channel <- data

	return nil
}

func (conveyer *Conveyer) Recv(output string) (string, error) {
	channel, ok := conveyer.channels[output] //nolint:varnamelen
	if !ok {
		return "", errRecvChanNotFound
	}

	data, ok := <-channel
	if !ok {
		return "undefined", nil
	}

	return data, nil
}
