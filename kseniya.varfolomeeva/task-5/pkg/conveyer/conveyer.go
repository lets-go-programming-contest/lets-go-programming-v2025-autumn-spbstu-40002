package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

const undefinedValue = "undefined"

type conveyer interface {
	RegisterDecorator(handler func(context.Context, chan string, chan string) error, input, output string)
	RegisterMultiplexer(handler func(context.Context, []chan string, chan string) error, inputs []string, output string)
	RegisterSeparator(handler func(context.Context, chan string, []chan string) error, input string, outputs []string)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

var _ conveyer = (*pipeline)(nil)

type workerFunc func(context.Context) error

type pipeline struct {
	bufferSize int
	channels   map[string]chan string
	workers    []workerFunc
	rwMutex    sync.RWMutex
}

func New(size int) *pipeline {
	return &pipeline{
		bufferSize: size,
		channels:   make(map[string]chan string),
		workers:    make([]workerFunc, 0),
		rwMutex:    sync.RWMutex{},
	}
}

func (p *pipeline) getOrCreateChannel(name string) chan string {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	if ch, ok := p.channels[name]; ok {
		return ch
	}

	newChan := make(chan string, p.bufferSize)
	p.channels[name] = newChan

	return newChan
}

func (p *pipeline) getChannel(name string) (chan string, error) {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	if ch, ok := p.channels[name]; ok {
		return ch, nil
	}

	return nil, ErrChanNotFound
}

func (p *pipeline) RegisterDecorator(
	handler func(context.Context, chan string, chan string) error,
	input, output string,
) {
	inChan := p.getOrCreateChannel(input)
	outChan := p.getOrCreateChannel(output)

	worker := func(ctx context.Context) error {
		return handler(ctx, inChan, outChan)
	}

	p.addWorker(worker)
}

func (p *pipeline) RegisterMultiplexer(
	handler func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	inChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChans[i] = p.getOrCreateChannel(name)
	}

	outChan := p.getOrCreateChannel(output)

	worker := func(ctx context.Context) error {
		return handler(ctx, inChans, outChan)
	}

	p.addWorker(worker)
}

func (p *pipeline) RegisterSeparator(
	handler func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	inChan := p.getOrCreateChannel(input)

	outChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChans[i] = p.getOrCreateChannel(name)
	}

	worker := func(ctx context.Context) error {
		return handler(ctx, inChan, outChans)
	}

	p.addWorker(worker)
}

func (p *pipeline) Send(input string, data string) error {
	ch, err := p.getChannel(input)
	if err != nil {
		return err
	}

	ch <- data

	return nil
}

func (p *pipeline) Recv(output string) (string, error) {
	ch, err := p.getChannel(output)
	if err != nil {
		return "", err
	}

	data, ok := <-ch
	if !ok {
		return undefinedValue, nil
	}

	return data, nil
}

func (p *pipeline) closeChannels() {
	for _, ch := range p.channels {
		close(ch)
	}
}

func (p *pipeline) Run(ctx context.Context) error {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	group, ctx := errgroup.WithContext(ctx)

	for _, w := range p.workers {
		workerFunc := w

		group.Go(func() error {
			return workerFunc(ctx)
		})
	}

	err := group.Wait()

	p.closeChannels()

	if err != nil {
		return fmt.Errorf("pipeline run failed: %w", err)
	}

	return nil
}

func (p *pipeline) addWorker(worker workerFunc) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	p.workers = append(p.workers, worker)
}
