package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type DecoratorFunc func(context.Context, chan string, chan string) error
type MultiplexerFunc func(context.Context, []chan string, chan string) error
type SeparatorFunc func(context.Context, chan string, []chan string) error

type Conveyer interface {
	RegisterDecorator(decorator DecoratorFunc, inputID string, outputID string) error
	RegisterMultiplexer(multiplexer MultiplexerFunc, inputIDs []string, outputID string) error
	RegisterSeparator(separator SeparatorFunc, inputID string, outputIDs []string) error
	Send(chanID string, value string) error
	Recv(chanID string) (string, error)
	Run(ctx context.Context) error
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

func (c *conv) ensureChan(chanID string) chan string {
	c.mu.Lock()
	existing, found := c.chans[chanID]
	if !found {
		existing = make(chan string, c.size)
		c.chans[chanID] = existing
	}
	c.mu.Unlock()

	return existing
}

func (c *conv) getChan(chanID string) (chan string, bool) {
	c.mu.RLock()
	existing, found := c.chans[chanID]
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

	value, more := <-ch
	if !more {
		return "", fmt.Errorf("%w: %s", ErrChannelClosed, chanID)
	}

	return value, nil
}

func (c *conv) RegisterDecorator(decorator DecoratorFunc, inputID string, outputID string) error {
	inputCh := c.ensureChan(inputID)
	outputCh := c.ensureChan(outputID)

	c.runners = append(c.runners, func(ctx context.Context) error {
		defer close(outputCh)

		return decorator(ctx, inputCh, outputCh)
	})

	return nil
}

func (c *conv) RegisterMultiplexer(multiplexer MultiplexerFunc, inputIDs []string, outputID string) error {
	inputs := make([]chan string, 0, len(inputIDs))
	for _, id := range inputIDs {
		inputs = append(inputs, c.ensureChan(id))
	}

	outputCh := c.ensureChan(outputID)

	c.runners = append(c.runners, func(ctx context.Context) error {
		defer close(outputCh)

		return multiplexer(ctx, inputs, outputCh)
	})

	return nil
}

func (c *conv) RegisterSeparator(separator SeparatorFunc, inputID string, outputIDs []string) error {
	inputCh := c.ensureChan(inputID)

	rawOutputs := make([]chan string, 0, len(outputIDs))
	outputs := make([]chan string, 0, len(outputIDs))

	for _, id := range outputIDs {
		ch := c.ensureChan(id)
		rawOutputs = append(rawOutputs, ch)
		outputs = append(outputs, ch)
	}

	c.runners = append(c.runners, func(ctx context.Context) error {
		for _, ch := range rawOutputs {
			defer close(ch)
		}

		return separator(ctx, inputCh, outputs)
	})

	return nil
}

func (c *conv) runAll(ctx context.Context) (<-chan struct{}, <-chan error) {
	var waitGroup sync.WaitGroup

	waitGroup.Add(len(c.runners))

	errOnce := make(chan error, 1)

	for _, r := range c.runners {
		runner := r

		go func() {
			defer waitGroup.Done()

			if err := runner(ctx); err != nil {
				select {
				case errOnce <- err:
				default:
				}
			}
		}()
	}

	doneCh := make(chan struct{})

	go func() {
		waitGroup.Wait()
		close(doneCh)
	}()

	return doneCh, errOnce
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
