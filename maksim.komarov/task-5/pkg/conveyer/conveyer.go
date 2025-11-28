package conveyer

import (
	"context"
	"errors"
	"sync"
)

const chanNotFoundMsg = "chan not found"

type handlerFunc func(ctx context.Context) error

type stringConveyer struct {
	size int

	mu       sync.RWMutex
	chans    map[string]chan string
	handlers []handlerFunc
}

func New(size int) *stringConveyer {
	if size < 0 {
		size = 0
	}
	return &stringConveyer{
		size:  size,
		chans: make(map[string]chan string),
	}
}

func (c *stringConveyer) ensureChan(id string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.chans[id]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.chans[id] = ch
	return ch
}

func (c *stringConveyer) getChan(id string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, ok := c.chans[id]
	return ch, ok
}

func (c *stringConveyer) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inCh := c.ensureChan(input)
	outCh := c.ensureChan(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	})
}

func (c *stringConveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inChans := make([]chan string, 0, len(inputs))
	for _, id := range inputs {
		inChans = append(inChans, c.ensureChan(id))
	}
	outCh := c.ensureChan(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inChans, outCh)
	})
}

func (c *stringConveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inCh := c.ensureChan(input)
	outChans := make([]chan string, 0, len(outputs))
	for _, id := range outputs {
		outChans = append(outChans, c.ensureChan(id))
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inCh, outChans)
	})
}

func (c *stringConveyer) Run(ctx context.Context) error {
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, h := range c.handlers {
		wg.Add(1)
		go func(h handlerFunc) {
			defer wg.Done()
			if err := h(runCtx); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				select {
				case errCh <- err:
				default:
				}
			}
		}(h)
	}

	var runErr error

	select {
	case <-ctx.Done():
		runErr = nil
	case err := <-errCh:
		runErr = err
	}

	cancel()
	wg.Wait()

	c.mu.Lock()
	for _, ch := range c.chans {
		close(ch)
	}
	c.mu.Unlock()

	return runErr
}

func (c *stringConveyer) Send(input string, data string) error {
	ch, ok := c.getChan(input)
	if !ok {
		return errors.New(chanNotFoundMsg)
	}
	ch <- data
	return nil
}

func (c *stringConveyer) Recv(output string) (string, error) {
	ch, ok := c.getChan(output)
	if !ok {
		return "", errors.New(chanNotFoundMsg)
	}

	v, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return v, nil
}
