package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type (
	DecoratorFunc   func(context.Context, chan string, chan string) error
	MultiplexerFunc func(context.Context, []chan string, chan string) error
	SeparatorFunc   func(context.Context, chan string, []chan string) error
)

type Conveyer interface {
	RegisterDecorator(DecoratorFunc, string, string) error
	RegisterMultiplexer(MultiplexerFunc, []string, string) error
	RegisterSeparator(SeparatorFunc, string, []string) error
	Send(string, string) error
	Recv(string) (string, error)
	Run(context.Context) error
}

var (
	ErrAlreadyRunning  = errors.New("already running")
	ErrChannelNotFound = errors.New("channel not found")
	ErrSend            = errors.New("send")
	ErrRecv            = errors.New("recv")
	ErrChannelClosed   = errors.New("channel closed")
)

type conv struct {
	size int

	mu      sync.RWMutex
	startMu sync.Mutex

	chans   map[string]chan string
	runners []func(context.Context) error
	started bool
}

func New(size int) *conv {
	if size < 0 {
		size = 0
	}

	return &conv{
		size:    size,
		mu:      sync.RWMutex{},
		startMu: sync.Mutex{},
		chans:   make(map[string]chan string),
		runners: make([]func(context.Context) error, 0),
		started: false,
	}
}

func (c *conv) ensureChan(id string) chan string {
	c.mu.Lock()
	existing, found := c.chans[id]

	if !found {
		existing = make(chan string, c.size)
		c.chans[id] = existing
	}

	c.mu.Unlock()

	return existing
}

func (c *conv) getChan(id string) (chan string, bool) {
	c.mu.RLock()
	existing, found := c.chans[id]
	c.mu.RUnlock()

	return existing, found
}

func (c *conv) Send(chanID string, value string) error {
	ch, ok := c.getChan(chanID)
	if !ok {
		return fmt.Errorf("%w: %w: %s", ErrChannelNotFound, ErrSend, chanID)
	}

	ch <- value

	return nil
}

func (c *conv) Recv(chanID string) (string, error) {
	ch, ok := c.getChan(chanID)
	if !ok {
		return "", fmt.Errorf("%w: %w: %s", ErrChannelNotFound, ErrRecv, chanID)
	}

	v, more := <-ch
	if !more {
		return "", fmt.Errorf("%w: %s", ErrChannelClosed, chanID)
	}

	return v, nil
}

func (c *conv) RegisterDecorator(fn DecoratorFunc, inputID string, outputID string) error {
	in := c.ensureChan(inputID)
	out := c.ensureChan(outputID)

	c.runners = append(c.runners, func(ctx context.Context) error {
		defer close(out)

		return fn(ctx, in, out)
	})

	return nil
}

func (c *conv) RegisterMultiplexer(fn MultiplexerFunc, inputIDs []string, outputID string) error {
	inputs := make([]chan string, 0, len(inputIDs))
	for _, id := range inputIDs {
		inputs = append(inputs, c.ensureChan(id))
	}

	out := c.ensureChan(outputID)

	c.runners = append(c.runners, func(ctx context.Context) error {
		defer close(out)

		return fn(ctx, inputs, out)
	})

	return nil
}

func (c *conv) RegisterSeparator(fn SeparatorFunc, inputID string, outputIDs []string) error {
	in := c.ensureChan(inputID)

	raw := make([]chan string, 0, len(outputIDs))
	outs := make([]chan string, 0, len(outputIDs))

	for _, id := range outputIDs {
		ch := c.ensureChan(id)
		raw = append(raw, ch)
		outs = append(outs, ch)
	}

	c.runners = append(c.runners, func(ctx context.Context) error {
		for _, ch := range raw {
			defer close(ch)
		}

		return fn(ctx, in, outs)
	})

	return nil
}

func (c *conv) runAll(ctx context.Context) (<-chan struct{}, <-chan error) {
	var wg sync.WaitGroup

	wg.Add(len(c.runners))

	errOnce := make(chan error, 1)

	for _, r := range c.runners {
		run := r

		go func() {
			defer wg.Done()

			if err := run(ctx); err != nil {
				select {
				case errOnce <- err:
				default:
				}
			}
		}()
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	return done, errOnce
}

func (c *conv) Run(ctx context.Context) error {
	c.startMu.Lock()
	if c.started {
		c.startMu.Unlock()

		return ErrAlreadyRunning
	}

	c.started = true
	c.startMu.Unlock()

	done, errOnce := c.runAll(ctx)

	select {
	case <-ctx.Done():
		<-done
	case <-done:
	}

	select {
	case err := <-errOnce:
		return err
	default:
		return nil
	}
}
