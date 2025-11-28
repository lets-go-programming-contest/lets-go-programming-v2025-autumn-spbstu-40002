package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

const undefined = "undefined"

type Pipeline struct {
	size     int
	channels map[string]chan string
	handlers []func(ctx context.Context) error
	mutex    sync.RWMutex
}

func New(size int) *Pipeline {
	return &Pipeline{
		size:     size,
		channels: make(map[string]chan string),
		mutex:    sync.RWMutex{},
		handlers: []func(ctx context.Context) error{},
	}
}

func (p *Pipeline) register(name string) chan string {
	if ch, ok := p.channels[name]; ok {
		return ch
	}

	ch := make(chan string, p.size)
	p.channels[name] = ch

	return ch
}

func (p *Pipeline) RegisterDecorator(
	callback func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	inCh := p.register(input)
	outCh := p.register(output)

	p.handlers = append(p.handlers, func(ctx context.Context) error {
		return callback(ctx, inCh, outCh)
	})
}

func (p *Pipeline) RegisterMultiplexer(
	callback func(
		ctx context.Context,
		inputs []chan string,
		output chan string,
	) error,
	inputs []string,
	output string,
) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	inChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChans[i] = p.register(name)
	}

	p.handlers = append(p.handlers, func(ctx context.Context) error {
		return callback(ctx, inChans, p.register(output))
	})
}

func (p *Pipeline) RegisterSeparator(
	callback func(
		ctx context.Context,
		input chan string,
		outputs []chan string,
	) error,
	input string,
	outputs []string,
) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	outChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChans[i] = p.register(name)
	}

	inCh := p.register(input)

	p.handlers = append(p.handlers, func(ctx context.Context) error {
		return callback(ctx, inCh, outChans)
	})
}

func (p *Pipeline) Run(ctx context.Context) error {
	defer func() {
		p.mutex.RLock()
		defer p.mutex.RUnlock()

		for _, ch := range p.channels {
			close(ch)
		}
	}()

	errgr, ctx := errgroup.WithContext(ctx)

	p.mutex.RLock()

	for _, handler := range p.handlers {
		errgr.Go(func() error {
			return handler(ctx)
		})
	}

	p.mutex.RUnlock()

	if err := errgr.Wait(); err != nil {
		return fmt.Errorf("run pipeline: %w", err)
	}

	return nil
}

func (p *Pipeline) Send(input string, data string) error {
	p.mutex.RLock()

	channel, ok := p.channels[input]

	p.mutex.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (p *Pipeline) Recv(output string) (string, error) {
	p.mutex.RLock()

	channel, ok := p.channels[output] //nolint:varnamelen

	p.mutex.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	data, ok := <-channel
	if !ok {
		return undefined, nil
	}

	return data, nil
}
