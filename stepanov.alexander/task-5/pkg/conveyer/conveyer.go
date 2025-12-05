package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

const Undefined = "undefined"

type conveyerImpl struct {
	channels map[string]chan string
	handlers []handler
	chanSize int
	mu       sync.RWMutex
	errCh    chan error
	wg       sync.WaitGroup
	running  bool
}

type handler func(ctx context.Context) error

func New(size int) *conveyerImpl {
	if size < 0 {
		size = 0
	}
	return &conveyerImpl{
		channels: make(map[string]chan string),
		chanSize: size,
		errCh:    make(chan error, 1),
	}
}

func (c *conveyerImpl) ensureChannel(name string) chan string {
	c.mu.RLock()
	ch, ok := c.channels[name]
	c.mu.RUnlock()

	if !ok {
		c.mu.Lock()
		defer c.mu.Unlock()
		if ch, ok = c.channels[name]; !ok {
			ch = make(chan string, c.chanSize)
			c.channels[name] = ch
		}
	}
	return ch
}

func (c *conveyerImpl) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	inputID string,
	outputID string,
) {
	inputCh := c.ensureChannel(inputID)
	outputCh := c.ensureChannel(outputID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputCh, outputCh)
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputsID []string,
	outputID string,
) {
	inputChannels := make([]chan string, len(inputsID))
	for i, id := range inputsID {
		inputChannels[i] = c.ensureChannel(id)
	}
	outputCh := c.ensureChannel(outputID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChannels, outputCh)
	})
}

func (c *conveyerImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	inputID string,
	outputsID []string,
) {
	inputCh := c.ensureChannel(inputID)
	outputChannels := make([]chan string, len(outputsID))
	for i, id := range outputsID {
		outputChannels[i] = c.ensureChannel(id)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputCh, outputChannels)
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return errors.New("conveyer already running")
	}
	c.running = true
	c.mu.Unlock()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, h := range c.handlers {
		c.wg.Add(1)
		go func(h handler) {
			defer c.wg.Done()
			if err := h(ctx); err != nil && err != context.Canceled && err != context.DeadlineExceeded {
				select {
				case c.errCh <- err:
				case <-ctx.Done():
				}
			}
		}(h)
	}

	var finalErr error
	select {
	case <-ctx.Done():
		finalErr = ctx.Err()
	case err := <-c.errCh:
		finalErr = err
		cancel()
	}

	c.wg.Wait()

	c.mu.Lock()
	for id, ch := range c.channels {
		select {
		case _, ok := <-ch:
			if !ok {
				continue
			}
		default:
		}
		func() {
			defer func() {
				if r := recover(); r != nil {
					_ = r
				}
			}()
			close(ch)
		}()
		_ = id
	}
	c.running = false
	c.mu.Unlock()

	return finalErr
}

func (c *conveyerImpl) Send(inputID string, data string) error {
	c.mu.RLock()
	ch, ok := c.channels[inputID]
	c.mu.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	defer func() {
		if r := recover(); r != nil {
			_ = r
		}
	}()

	ch <- data
	return nil
}

func (c *conveyerImpl) Recv(outputID string) (string, error) {
	c.mu.RLock()
	ch, ok := c.channels[outputID]
	c.mu.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	data, open := <-ch
	if !open {
		return Undefined, nil
	}

	return data, nil
}
