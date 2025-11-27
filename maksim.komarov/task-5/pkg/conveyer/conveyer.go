package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type (
	DecoratorFunc   func(ctx context.Context, input chan string, output chan string) error
	MultiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error
	SeparatorFunc   func(ctx context.Context, input chan string, outputs []chan string) error
)

const Undefined = "undefined"

var (
	ErrAlreadyRunning = errors.New("conveyer already running")
	ErrChanNotFound   = errors.New("chan not found")
)

type Conveyer interface {
	RegisterDecorator(decorator DecoratorFunc, inputID string, outputID string) error
	RegisterMultiplexer(multiplexer MultiplexerFunc, inputIDs []string, outputID string) error
	RegisterSeparator(separator SeparatorFunc, inputID string, outputIDs []string) error
	Run(ctx context.Context) error
	Send(inputID string, data string) error
	Recv(outputID string) (string, error)
}

type runFunc func(ctx context.Context) error

type conv struct {
	bufferSize int
	chans      map[string]chan string
	runners    []runFunc
	started    bool
	mu         sync.RWMutex
	startMux   sync.Mutex
}

func New(size int) *conv {
	return &conv{
		bufferSize: size,
		chans:      make(map[string]chan string),
		runners:    make([]runFunc, 0),
		started:    false,
		mu:         sync.RWMutex{},
		startMux:   sync.Mutex{},
	}
}

func (c *conv) ensureChan(id string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	existing, found := c.chans[id]
	if !found {
		existing = make(chan string, c.bufferSize)
		c.chans[id] = existing
	}

	return existing
}

func (c *conv) getChan(id string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	existing, found := c.chans[id]

	return existing, found
}

func (c *conv) RegisterDecorator(fn DecoratorFunc, inputID string, outputID string) error {
	in := c.ensureChan(inputID)
	out := c.ensureChan(outputID)

	c.runners = append(c.runners, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})

	return nil
}

func (c *conv) RegisterMultiplexer(fn MultiplexerFunc, inputIDs []string, outputID string) error {
	inputs := make([]chan string, 0, len(inputIDs))
	for _, one := range inputIDs {
		inputs = append(inputs, c.ensureChan(one))
	}

	out := c.ensureChan(outputID)

	c.runners = append(c.runners, func(ctx context.Context) error {
		return fn(ctx, inputs, out)
	})

	return nil
}

func (c *conv) RegisterSeparator(fn SeparatorFunc, inputID string, outputIDs []string) error {
	in := c.ensureChan(inputID)

	outs := make([]chan string, 0, len(outputIDs))
	for _, one := range outputIDs {
		outs = append(outs, c.ensureChan(one))
	}

	c.runners = append(c.runners, func(ctx context.Context) error {
		return fn(ctx, in, outs)
	})

	return nil
}

func (c *conv) Run(ctx context.Context) error {
	c.startMux.Lock()
	if c.started {
		c.startMux.Unlock()

		return ErrAlreadyRunning
	}
	c.started = true
	c.startMux.Unlock()

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error, len(c.runners))

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(c.runners))

	for _, r := range c.runners {
		runner := r

		go func() {
			defer waitGroup.Done()

			if err := runner(runCtx); err != nil {
				errChan <- err
			}
		}()
	}

	var firstErr error

	select {
	case <-ctx.Done():
		firstErr = fmt.Errorf("context canceled: %w", ctx.Err())
	case err := <-errChan:
		firstErr = err
	}

	cancel()

	c.mu.Lock()
	for id, ch := range c.chans {
		close(ch)
		delete(c.chans, id)
	}
	c.mu.Unlock()

	waitGroup.Wait()

	return firstErr
}

func (c *conv) Send(inputID string, data string) error {
	ch, found := c.getChan(inputID)
	if !found {
		return fmt.Errorf("%w", ErrChanNotFound)
	}

	ch <- data

	return nil
}

func (c *conv) Recv(outputID string) (string, error) {
	ch, found := c.getChan(outputID)
	if !found {
		return "", fmt.Errorf("%w", ErrChanNotFound)
	}

	value, ok := <-ch
	if !ok {
		return Undefined, nil
	}

	return value, nil
}
