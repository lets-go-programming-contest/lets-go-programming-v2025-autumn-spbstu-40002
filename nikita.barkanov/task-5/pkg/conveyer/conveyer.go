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
	fn func(context.Context, chan string, chan string) error,
	input, output string,
) {
	inCh := conv.createOrGetChannel(input)
	outCh := conv.createOrGetChannel(output)

	runner := func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	}

	conv.rwmu.Lock()
	conv.handlers = append(conv.handlers, runner)
	conv.rwmu.Unlock()
}

func (conv *conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string, output string,
) {
	outCh := conv.createOrGetChannel(output)

	inputChs := make([]chan string, 0, len(inputs))
	for _, name := range inputs {
		ch := conv.createOrGetChannel(name)
		inputChs = append(inputChs, ch)
	}

	runner := func(ctx context.Context) error {
		return fn(ctx, inputChs, outCh)
	}

	conv.rwmu.Lock()
	conv.handlers = append(conv.handlers, runner)
	conv.rwmu.Unlock()
}

func (conv *conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string, outputs []string,
) {
	inCh := conv.createOrGetChannel(input)

	outputChs := make([]chan string, 0, len(outputs))
	for _, name := range outputs {
		ch := conv.createOrGetChannel(name)
		outputChs = append(outputChs, ch)
	}

	runner := func(ctx context.Context) error {
		return fn(ctx, inCh, outputChs)
	}

	conv.rwmu.Lock()
	conv.handlers = append(conv.handlers, runner)
	conv.rwmu.Unlock()
}

func (conv *conveyer) Run(ctx context.Context) error {
	if len(conv.handlers) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	conv.rwmu.RLock()
	handlers := make([]func(context.Context) error, len(conv.handlers))
	copy(handlers, conv.handlers)
	conv.rwmu.RUnlock()

	for _, h := range handlers {
		wg.Add(1)
		handler := h
		go func() {
			defer wg.Done()
			if err := handler(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	select {
	case err, ok := <-errCh:
		if ok {
			return err
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func New(size int) Conveyer {
	return &conveyer{
		channels: make(map[string]chan string),
		size:     size,
	}
}
