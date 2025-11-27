package conveyer

import (
	"context"
	"errors"
	"sync"
)

type DecoratorFunc func(context.Context, <-chan string, chan<- string) error
type MultiplexerFunc func(context.Context, []<-chan string, chan<- string) error
type SeparatorFunc func(context.Context, <-chan string, []chan<- string) error

type Conveyer interface {
	RegisterDecorator(DecoratorFunc, string, string) error
	RegisterMultiplexer(MultiplexerFunc, []string, string) error
	RegisterSeparator(SeparatorFunc, string, []string) error
	Send(string, string)
	Recv(string) (string, bool)
	Run(context.Context) error
}

type conv struct {
	size    int
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
		chans:   make(map[string]chan string),
		runners: make([]func(context.Context) error, 0),
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

func (c *conv) Send(chanID, value string) {
	ch := c.ensureChan(chanID)

	ch <- value
}

func (c *conv) Recv(chanID string) (string, bool) {
	ch, found := c.getChan(chanID)
	if !found {
		return "", false
	}

	value, ok := <-ch
	if !ok {
		return "", false
	}

	return value, true
}

func (c *conv) RegisterDecorator(handler DecoratorFunc, inputID, outputID string) error {
	inputCh := c.ensureChan(inputID)
	outputCh := c.ensureChan(outputID)

	c.runners = append(c.runners, func(ctx context.Context) error {
		defer close(outputCh)

		return handler(ctx, inputCh, outputCh)
	})

	return nil
}

func (c *conv) RegisterMultiplexer(handler MultiplexerFunc, inputIDs []string, outputID string) error {
	inputs := make([]<-chan string, 0, len(inputIDs))
	for _, id := range inputIDs {
		inputs = append(inputs, c.ensureChan(id))
	}
	outputCh := c.ensureChan(outputID)

	c.runners = append(c.runners, func(ctx context.Context) error {
		defer close(outputCh)

		return handler(ctx, inputs, outputCh)
	})

	return nil
}

func (c *conv) RegisterSeparator(handler SeparatorFunc, inputID string, outputIDs []string) error {
	inputCh := c.ensureChan(inputID)

	raw := make([]chan string, 0, len(outputIDs))
	outs := make([]chan<- string, 0, len(outputIDs))
	for _, id := range outputIDs {
		ch := c.ensureChan(id)
		raw = append(raw, ch)
		outs = append(outs, ch)
	}

	c.runners = append(c.runners, func(ctx context.Context) error {
		for _, ch := range raw {
			defer close(ch)
		}

		return handler(ctx, inputCh, outs)
	})

	return nil
}

func (c *conv) Run(ctx context.Context) error {
	c.startMu.Lock()
	if c.started {
		c.startMu.Unlock()

		return errors.New("already running")
	}
	c.started = true
	c.startMu.Unlock()

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

	done := make(chan struct{})
	go func() {
		waitGroup.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		<-done

		select {
		case err := <-errOnce:
			if err != nil {
				return err
			}
		default:
		}

		return nil

	case <-done:
		select {
		case err := <-errOnce:
			if err != nil {
				return err
			}
		default:
		}

		return nil
	}
}
