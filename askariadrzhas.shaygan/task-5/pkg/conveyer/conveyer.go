package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Pipeline interface {
	AddTransformer(
		processor func(
			ctx context.Context,
			inCh chan string,
			outCh chan string,
		) error,
		source string,
		dest string,
	)

	AddMerger(
		processor func(
			ctx context.Context,
			sources []chan string,
			dest chan string,
		) error,
		sources []string,
		dest string,
	)

	AddSplitter(
		processor func(
			ctx context.Context,
			source chan string,
			dests []chan string,
		) error,
		source string,
		dests []string,
	)

	Start(ctx context.Context) error
	Push(source string, value string) error
	Pull(dest string) (string, error)
}

type PipelineImpl struct {
	chans      map[string]chan string
	bufferSize int
	processors []func(ctx context.Context) error
	accessLock sync.RWMutex
}

func Create(size int) PipelineImpl {
	return PipelineImpl{
		chans:      make(map[string]chan string),
		bufferSize: size,
		processors: make([]func(ctx context.Context) error, 0),
		accessLock: sync.RWMutex{},
	}
}

func (p *PipelineImpl) ensureChan(name string) {
	p.accessLock.Lock()
	defer p.accessLock.Unlock()

	if _, exists := p.chans[name]; !exists {
		p.chans[name] = make(chan string, p.bufferSize)
	}
}

func (p *PipelineImpl) fetchChan(name string) (chan string, bool) {
	if ch, ok := p.chans[name]; ok {
		return ch, true
	}
	return nil, false
}

func (p *PipelineImpl) shutdownChans() {
	p.accessLock.Lock()
	defer p.accessLock.Unlock()

	for _, ch := range p.chans {
		close(ch)
	}
}

func (p *PipelineImpl) registerProcessor(processor func(ctx context.Context) error) {
	p.processors = append(p.processors, processor)
}

func (p *PipelineImpl) AddTransformer(
	processor func(
		ctx context.Context,
		inCh chan string,
		outCh chan string,
	) error,
	source string, dest string,
) {
	p.ensureChan(source)
	p.ensureChan(dest)

	p.registerProcessor(func(ctx context.Context) error {
		p.accessLock.RLock()
		defer p.accessLock.RUnlock()

		inCh, _ := p.fetchChan(source)
		outCh, _ := p.fetchChan(dest)

		return processor(ctx, inCh, outCh)
	})
}

func (p *PipelineImpl) AddMerger(
	processor func(
		ctx context.Context,
		sources []chan string,
		dest chan string,
	) error,
	sources []string, dest string,
) {
	for _, src := range sources {
		p.ensureChan(src)
	}
	p.ensureChan(dest)

	p.registerProcessor(func(ctx context.Context) error {
		p.accessLock.RLock()
		defer p.accessLock.RUnlock()

		sourceChans := make([]chan string, len(sources))
		for i, src := range sources {
			ch, _ := p.fetchChan(src)
			sourceChans[i] = ch
		}

		destChan, _ := p.fetchChan(dest)
		return processor(ctx, sourceChans, destChan)
	})
}

func (p *PipelineImpl) AddSplitter(
	processor func(
		ctx context.Context,
		source chan string,
		dests []chan string,
	) error,
	source string, dests []string,
) {
	p.ensureChan(source)
	for _, dest := range dests {
		p.ensureChan(dest)
	}

	p.registerProcessor(func(ctx context.Context) error {
		p.accessLock.RLock()
		defer p.accessLock.RUnlock()

		sourceChan, _ := p.fetchChan(source)
		destChans := make([]chan string, len(dests))
		for i, dest := range dests {
			ch, _ := p.fetchChan(dest)
			destChans[i] = ch
		}

		return processor(ctx, sourceChan, destChans)
	})
}

func (p *PipelineImpl) Start(ctx context.Context) error {
	defer p.shutdownChans()

	group, ctxWithCancel := errgroup.WithContext(ctx)

	p.accessLock.RLock()
	for _, processor := range p.processors {
		proc := processor
		group.Go(func() error {
			return proc(ctxWithCancel)
		})
	}
	p.accessLock.RUnlock()

	if err := group.Wait(); err != nil {
		return err
	}
	return nil
}

func (p *PipelineImpl) Push(source string, value string) error {
	p.accessLock.RLock()
	ch, found := p.fetchChan(source)
	p.accessLock.RUnlock()

	if !found {
		return ErrChanMissing
	}

	ch <- value
	return nil
}

func (p *PipelineImpl) Pull(dest string) (string, error) {
	p.accessLock.RLock()
	ch, found := p.fetchChan(dest)
	p.accessLock.RUnlock()

	if !found {
		return "", ErrChanMissing
	}

	if data, ok := <-ch; ok {
		return data, nil
	}
	return Undefined, nil
}

var (
	ErrChanMissing = errors.New("channel does not exist")
	Undefined      = "undefined"
)
