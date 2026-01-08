package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

const UndefinedValue = "undefined"

type ConveyerInterface interface {
	RegisterDecorator(handler func(ctx context.Context, input chan string, output chan string) error, input, output string)
	RegisterMultiplexer(handler func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string)
	RegisterSeparator(handler func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type pipeline struct {
	bufferSize int
	channels   map[string]chan string
	workers    []func(ctx context.Context) error
	rwMutex    sync.RWMutex
}

func New(size int) *pipeline {
	return &pipeline{
		bufferSize: size,
		channels:   make(map[string]chan string),
		workers:    make([]func(ctx context.Context) error, 0),
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
	handler func(ctx context.Context, input chan string, output chan string) error,
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
	handler func(ctx context.Context, inputs []chan string, output chan string) error,
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
	handler func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inChan := p.getOrCreateChannel(input)

	outList := make([]chan string, len(outputs))
	for i, name := range outputs {
		outList[i] = p.getOrCreateChannel(name)
	}

	worker := func(ctx context.Context) error {
		return handler(ctx, inChan, outList)
	}

	p.addWorker(worker)
}

func (p *pipeline) Send(input string, data string) error {
	ch, err := p.getChannel(input)
	if err != nil {
		return err
	}

	select {
	case ch <- data:
		return nil
	default:
		return fmt.Errorf("channel %s is full", input)
	}
}

func (p *pipeline) Recv(output string) (string, error) {
	ch, err := p.getChannel(output)
	if err != nil {
		return "", err
	}

	select {
	case value := <-ch:
		return value, nil
	default:
		return UndefinedValue, nil
	}
}

func (p *pipeline) closeChannels() {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	for _, ch := range p.channels {
		close(ch)
	}
}

func (p *pipeline) Run(ctx context.Context) error {
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

func (p *pipeline) addWorker(worker func(ctx context.Context) error) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	p.workers = append(p.workers, worker)
}
