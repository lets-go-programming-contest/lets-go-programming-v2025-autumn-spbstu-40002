package conveyer

import (
	"context"
	"errors"
	"sync"
)

var errChanNotFound = errors.New("chan not found")

type Conveyer interface {
	RegisterDecorator(
		fn func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		fn func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		fn func(ctx context.Context, input chan string, outputs []chan string) error,
		input string,
		outputs []string,
	)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type conveyer struct {
	rwmu     sync.RWMutex
	channels map[string]chan string
	handlers []func(ctx context.Context) error
	size     int
}

func (conv *conveyer) getChannel(name string) (chan string, bool) {
	conv.rwmu.RLock()
	defer conv.rwmu.RUnlock()
	channel, isOkey := conv.channels[name]

	return channel, isOkey
}

func (conv *conveyer) Send(input string, data string) error {
	channel, exists := conv.getChannel(input)
	if !exists {
		return errChanNotFound
	}

	channel <- data

	return nil
}

func (conv *conveyer) Recv(output string) (string, error) {
	channel, exists := conv.getChannel(output)
	if !exists {
		return "", errChanNotFound
	}

	data, isOkey := <-channel
	if !isOkey {
		return "undefined", nil
	}

	return data, nil
}

func (conv *conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	input, output string,
) {

}

func (conv *conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string, output string,
) {

}

func (conv *conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string, outputs []string,
) {
}

func (conv *conveyer) Run(ctx context.Context) error {
	return nil
}

func New(size int) Conveyer {
	return &conveyer{
		channels: make(map[string]chan string),
		size:     size,
	}
}
