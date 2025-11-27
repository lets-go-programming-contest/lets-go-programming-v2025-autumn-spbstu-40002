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
	RegisterDecorator(fn DecoratorFunc, inputID string, outputID string) error
	RegisterMultiplexer(fn MultiplexerFunc, inputIDs []string, outputID string) error
	RegisterSeparator(fn SeparatorFunc, inputID string, outputIDs []string) error
	Send(chanID string, value string) error
	Recv(chanID string) (string, error)
	Run(ctx context.Context) error
}

var ErrAlreadyRunning = errors.New("already running")
var ErrChannelNotFound = errors.New("channel not found")

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
		return errors.Join(ErrChannelNotFound, errors.New("send: "+chanID))
	}

	ch <- value

	return nil
}

func (c *conv) Recv(chanID string) (string, error) {
	ch, ok := c.getChan(chanID)
	if !ok {
		return "", errors.Join(ErrChannelNotFound, errors.New("recv: "+chanID))
	}

	val, more := <-ch
	if !more {
		return "", errors.New("channel closed: " + chanID)
	}

	return val, nil
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
	ins := make([]<-chan string, 0, len(inputIDs))
	for _, id := range inputIDs {
		ins = append(ins, c.ensureChan(id))
	}

	out := c.ensureChan(outputID)

	c.runners = append(c.runners, func(ctx context.Context) error {
		defer close(out)

		return fn(ctx, ins, out)
	})

	return nil
}

func (c *conv) RegisterSeparator(fn SeparatorFunc, inputID string, outputIDs []string) error {
	in := c.ensureChan(inputID)

	rawOuts := make([]chan string, 0, len(outputIDs))
	outs := make([]chan<- string, 0, len(outputIDs))

	for _, id := range outputIDs {
		ch := c.ensureChan(id)
		rawOuts = append(rawOuts, ch)
		outs = append(outs, ch)
	}

	c.runners = append(c.runners, func(ctx context.Context) error {
		for _, ch := range rawOuts {
			defer close(ch)
		}

		return fn(ctx, in, outs)
	})

	return nil
}

func (c *conv) runAll(ctx context.Context) (done <-chan struct{}, errOnce <-chan error) {
	var wg sync.WaitGroup

	wg.Add(len(c.runners))

	errChan := make(chan error, 1)
	for _, r := range c.runners {
		task := r

		go func() {
			defer wg.Done()

			if err := task(ctx); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		}()
	}

	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	return doneCh, errChan
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
