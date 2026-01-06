package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	errChanNotFound = errors.New("chan not found")
	errWaitConveyer = errors.New("can't be decorated")
)

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

func (conv *conveyer) createOrGetChannel(name string) chan string {
	conv.rwmu.Lock()
	defer conv.rwmu.Unlock()

	if ch, ok := conv.channels[name]; ok {
		return ch
	}

	ch := make(chan string, conv.size)
	conv.channels[name] = ch
	return ch
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
	handlerFunc func(context.Context, chan string, chan string) error,
	input, output string,
) {
	inCh := conv.createOrGetChannel(input)
	outCh := conv.createOrGetChannel(output)

	runner := func(ctx context.Context) error {
		return handlerFunc(ctx, inCh, outCh)
	}

	conv.rwmu.Lock()
	conv.handlers = append(conv.handlers, runner)
	conv.rwmu.Unlock()
}

func (conv *conveyer) RegisterMultiplexer(
	handlerFunc func(context.Context, []chan string, chan string) error,
	inputs []string, output string,
) {
	outCh := conv.createOrGetChannel(output)

	inputChs := make([]chan string, 0, len(inputs))

	for _, name := range inputs {
		ch := conv.createOrGetChannel(name)
		inputChs = append(inputChs, ch)
	}

	runner := func(ctx context.Context) error {
		return handlerFunc(ctx, inputChs, outCh)
	}

	conv.rwmu.Lock()
	conv.handlers = append(conv.handlers, runner)
	conv.rwmu.Unlock()
}

func (conv *conveyer) RegisterSeparator(
	handlerFunc func(context.Context, chan string, []chan string) error,
	input string, outputs []string,
) {
	inCh := conv.createOrGetChannel(input)

	outputChs := make([]chan string, 0, len(outputs))

	for _, name := range outputs {
		ch := conv.createOrGetChannel(name)
		outputChs = append(outputChs, ch)
	}

	runner := func(ctx context.Context) error {
		return handlerFunc(ctx, inCh, outputChs)
	}

	conv.rwmu.Lock()
	conv.handlers = append(conv.handlers, runner)
	conv.rwmu.Unlock()
}

func (conv *conveyer) Run(ctx context.Context) error {
	if len(conv.handlers) == 0 {
		return nil
	}

	conv.rwmu.RLock()
	handlers := make([]func(context.Context) error, len(conv.handlers))
	copy(handlers, conv.handlers)
	conv.rwmu.RUnlock()

	group, gctx := errgroup.WithContext(ctx)

	for _, h := range handlers {
		handler := h

		group.Go(func() error {
			return handler(gctx)
		})
	}

	if err := group.Wait(); err != nil {
		return errWaitConveyer
	}

	return nil
}

func New(size int) *conveyer {
	return &conveyer{
		channels: make(map[string]chan string),
		size:     size,
		rwmu:     sync.RWMutex{},
		handlers: []func(ctx context.Context) error{},
	}
}
