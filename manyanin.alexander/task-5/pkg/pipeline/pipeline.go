package pipeline

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var chanNotFound = fmt.Errorf("channel not found")

const notDefined = "not_defined"

type Pipeline struct {
	bufferSize int
	channels   map[string]chan string
	processors []func(ctx context.Context) error
}

func (pl *Pipeline) getOrCreateChannel(name string) chan string {
	if ch, ok := pl.channels[name]; ok {
		return ch
	}

	ch := make(chan string, pl.bufferSize)
	pl.channels[name] = ch

	return ch
}

func New(size int) *Pipeline {
	if size < 0 {
		size = 0
	}

	return &Pipeline{
		bufferSize: size,
		channels:   make(map[string]chan string),
		processors: make([]func(ctx context.Context) error, 0),
	}
}

func (pl *Pipeline) AddTransformer(
	processor func(
		ctx context.Context,
		inCh chan string,
		outCh chan string,
	) error,
	inName string,
	outName string,
) {
	inChan := pl.getOrCreateChannel(inName)
	outChan := pl.getOrCreateChannel(outName)
	pl.processors = append(pl.processors, func(ctx context.Context) error {
		return processor(ctx, inChan, outChan)
	})
}

func (pl *Pipeline) AddMerger(
	processor func(
		ctx context.Context,
		inChs []chan string,
		outCh chan string,
	) error,
	inNames []string,
	outName string,
) {
	inChannels := make([]chan string, len(inNames))

	for i, name := range inNames {
		inChannels[i] = pl.getOrCreateChannel(name)
	}

	outChan := pl.getOrCreateChannel(outName)
	pl.processors = append(pl.processors, func(ctx context.Context) error {
		return processor(ctx, inChannels, outChan)
	})
}

func (pl *Pipeline) AddSplitter(
	processor func(
		ctx context.Context,
		inCh chan string,
		outChs []chan string,
	) error,
	inName string,
	outNames []string,
) {
	inChan := pl.getOrCreateChannel(inName)
	outChannels := make([]chan string, len(outNames))

	for i, name := range outNames {
		outChannels[i] = pl.getOrCreateChannel(name)
	}

	pl.processors = append(pl.processors, func(ctx context.Context) error {
		return processor(ctx, inChan, outChannels)
	})
}

func (pl *Pipeline) Run(ctx context.Context) error {
	group, ctxWithCancel := errgroup.WithContext(ctx)

	for _, proc := range pl.processors {
		currentProc := proc
		group.Go(func() error {
			return currentProc(ctxWithCancel)
		})
	}

	err := group.Wait()
	if err != nil {
		return fmt.Errorf("pipeline execution: %w", err)
	}

	return nil
}

func (pl *Pipeline) SendData(channelName string, data string) error {
	ch, exists := pl.channels[channelName]
	if !exists {
		return chanNotFound
	}

	ch <- data
	return nil
}

func (pl *Pipeline) ReceiveData(channelName string) (string, error) {
	ch, exists := pl.channels[channelName]
	if !exists {
		return "", chanNotFound
	}

	val, ok := <-ch
	if !ok {
		return notDefined, nil
	}

	return val, nil
}

func (pl *Pipeline) CloseAllChannels() {
	var wg sync.WaitGroup
	for name, ch := range pl.channels {
		wg.Add(1)
		go func(n string, c chan string) {
			defer wg.Done()
			close(c)
		}(name, ch)
	}
	wg.Wait()
}
