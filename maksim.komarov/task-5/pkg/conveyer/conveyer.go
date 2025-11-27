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
	ErrNotStarted     = errors.New("conveyer not started")
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

func (c *conv) ensureChan(chanID string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()
	existing, found := c.chans[chanID]
	if !found {
		existing = make(chan string, c.bufferSize)
		c.chans[chanID] = existing
	}
	return existing
}

func (c *conv) getChan(chanID string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	existing, found := c.chans[chanID]
	return existing, found
}

func (c *conv) RegisterDecorator(decorator DecoratorFunc, inputID string, outputID string) error {
	inputChan := c.ensureChan(inputID)
	outputChan := c.ensureChan(outputID)
	c.runners = append(c.runners, func(ctx context.Context) error {
		return decorator(ctx, inputChan, outputChan)
	})
	return nil
}

func (c *conv) RegisterMultiplexer(multiplexer MultiplexerFunc, inputIDs []string, outputID string) error {
	inputs := make([]chan string, 0, len(inputIDs))
	for _, oneID := range inputIDs {
		inputs = append(inputs, c.ensureChan(oneID))
	}
	outputChan := c.ensureChan(outputID)
	c.runners = append(c.runners, func(ctx context.Context) error {
		return multiplexer(ctx, inputs, outputChan)
	})
	return nil
}

func (c *conv) RegisterSeparator(separator SeparatorFunc, inputID string, outputIDs []string) error {
	inputChan := c.ensureChan(inputID)
	outputs := make([]chan string, 0, len(outputIDs))
	for _, oneID := range outputIDs {
		outputs = append(outputs, c.ensureChan(oneID))
	}
	c.runners = append(c.runners, func(ctx context.Context) error {
		return separator(ctx, inputChan, outputs)
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

	for _, runner := range c.runners {
		r := runner
		go func() {
			defer waitGroup.Done()
			if err := r(runCtx); err != nil {
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
	targetChan, found := c.getChan(inputID)
	if !found {
		return fmt.Errorf("%w", ErrChanNotFound)
	}
	targetChan <- data
	return nil
}

func (c *conv) Recv(outputID string) (string, error) {
	outputChan, found := c.getChan(outputID)
	if !found {
		return "", fmt.Errorf("%w", ErrChanNotFound)
	}
	value, ok := <-outputChan
	if !ok {
		return Undefined, nil
	}
	return value, nil
}
