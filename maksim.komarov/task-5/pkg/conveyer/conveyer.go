package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type DecoratorFunc func(ctx context.Context, input chan string, output chan string) error
type MultiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error
type SeparatorFunc func(ctx context.Context, input chan string, outputs []chan string) error

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
	ErrChannelNotFound = errors.New("chan not found")
)

type conv struct {
	mu      sync.Mutex
	startMu sync.Mutex
	started bool
	chans   map[string]chan string
	runners []func(ctx context.Context) error
}

func New(size int) *conv {
	return &conv{
		mu:      sync.Mutex{},
		startMu: sync.Mutex{},
		started: false,
		chans:   make(map[string]chan string, size),
		runners: make([]func(context.Context) error, 0),
	}
}

func (c *conv) ensureChan(chanID string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if existing, found := c.chans[chanID]; found {
		return existing
	}

	ch := make(chan string)
	c.chans[chanID] = ch

	return ch
}

func (c *conv) getChan(chanID string) (chan string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch, found := c.chans[chanID]

	return ch, found
}

func (c *conv) RegisterDecorator(decorator DecoratorFunc, inputID string, outputID string) error {
	inputChan := c.ensureChan(inputID)
	outputChan := c.ensureChan(outputID)

	c.runners = append(c.runners, func(ctx context.Context) error {
		defer close(outputChan)

		return decorator(ctx, inputChan, outputChan)
	})

	return nil
}

func (c *conv) RegisterMultiplexer(multiplexer MultiplexerFunc, inputIDs []string, outputID string) error {
	inputs := make([]chan string, 0, len(inputIDs))
	for _, id := range inputIDs {
		inputs = append(inputs, c.ensureChan(id))
	}

	outputChan := c.ensureChan(outputID)

	c.runners = append(c.runners, func(ctx context.Context) error {
		defer close(outputChan)

		return multiplexer(ctx, inputs, outputChan)
	})

	return nil
}

func (c *conv) RegisterSeparator(separator SeparatorFunc, inputID string, outputIDs []string) error {
	inputChan := c.ensureChan(inputID)
	outputs := make([]chan string, 0, len(outputIDs))
	for _, id := range outputIDs {
		outputs = append(outputs, c.ensureChan(id))
	}

	c.runners = append(c.runners, func(ctx context.Context) error {
		defer func() {
			for _, ch := range outputs {
				close(ch)
			}
		}()

		return separator(ctx, inputChan, outputs)
	})

	return nil
}

func (c *conv) Send(chanID string, value string) error {
	ch, ok := c.getChan(chanID)
	if !ok {
		return fmt.Errorf("%w: send: %s", ErrChannelNotFound, chanID)
	}

	ch <- value

	return nil
}

func (c *conv) Recv(chanID string) (string, error) {
	ch, ok := c.getChan(chanID)
	if !ok {
		return "", fmt.Errorf("%w: recv: %s", ErrChannelNotFound, chanID)
	}

	val, opened := <-ch
	if !opened {
		return "", nil
	}

	return val, nil
}

func (c *conv) runAll(ctx context.Context) (<-chan struct{}, <-chan error) {
	wg := &sync.WaitGroup{}
	wg.Add(len(c.runners))

	errChan := make(chan error, 1)

	for _, r := range c.runners {
		runFunc := r

		go func() {
			defer wg.Done()

			if err := runFunc(ctx); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		}()
	}

	doneChan := make(chan struct{})

	go func() {
		wg.Wait()
		close(doneChan)
	}()

	return doneChan, errChan
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
		return nil
	case err := <-errOnce:
		return err
	case <-done:
		return nil
	}
}
