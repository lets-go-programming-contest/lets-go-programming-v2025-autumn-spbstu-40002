package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	eg "golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound           = errors.New("chan not found")
	ErrConveyerNotRunning     = errors.New("conveyer is not running")
	ErrConveyerAlreadyRunning = errors.New("conveyer is already running")
	ErrChannelFull            = errors.New("channel is full")
	ErrHandlerAfterRun        = errors.New("cannot register handlers after Run() has been called")
)

type impl struct {
	mu sync.RWMutex

	channels    map[string]chan string
	channelSize int

	decorators   []decoratorHandler
	multiplexers []multiplexerHandler
	separators   []separatorHandler

	ctx    context.Context
	cancel context.CancelFunc
	g      *eg.Group

	running bool
	stopped bool
}

type decoratorHandler struct {
	fn     func(context.Context, chan string, chan string) error
	input  string
	output string
}

type multiplexerHandler struct {
	fn     func(context.Context, []chan string, chan string) error
	inputs []string
	output string
}

type separatorHandler struct {
	fn      func(context.Context, chan string, []chan string) error
	input   string
	outputs []string
}

func (c *impl) getOrCreateChannel(id string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[id]; exists {
		return ch
	}

	ch := make(chan string, c.channelSize)
	c.channels[id] = ch
	return ch
}

func (c *impl) getChannel(id string) chan string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.channels[id]
}

func (c *impl) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		panic(ErrHandlerAfterRun)
	}

	c.decorators = append(c.decorators, decoratorHandler{
		fn:     fn,
		input:  input,
		output: output,
	})
}

func (c *impl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		panic(ErrHandlerAfterRun)
	}

	c.multiplexers = append(c.multiplexers, multiplexerHandler{
		fn:     fn,
		inputs: inputs,
		output: output,
	})
}

func (c *impl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		panic(ErrHandlerAfterRun)
	}

	c.separators = append(c.separators, separatorHandler{
		fn:      fn,
		input:   input,
		outputs: outputs,
	})
}

func (c *impl) Send(input string, data string) error {
	c.mu.RLock()
	running := c.running
	stopped := c.stopped
	c.mu.RUnlock()

	if !running || stopped {
		return ErrConveyerNotRunning
	}

	ch := c.getChannel(input)
	if ch == nil {
		return ErrChanNotFound
	}

	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	case ch <- data:
		return nil
	default:
		return fmt.Errorf("%s: %s", ErrChannelFull, input)
	}
}

func (c *impl) Recv(output string) (string, error) {
	c.mu.RLock()
	running := c.running
	stopped := c.stopped
	c.mu.RUnlock()

	if !running || stopped {
		return "", ErrConveyerNotRunning
	}

	ch := c.getChannel(output)
	if ch == nil {
		return "", ErrChanNotFound
	}

	select {
	case <-c.ctx.Done():
		return "", c.ctx.Err()
	case data, ok := <-ch:
		if !ok {
			return "undefined", nil
		}
		return data, nil
	}
}

func (c *impl) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return ErrConveyerAlreadyRunning
	}
	c.running = true
	c.mu.Unlock()

	c.ctx, c.cancel = context.WithCancel(ctx)
	c.g, c.ctx = eg.WithContext(c.ctx)

	c.createAllChannels()
	c.startHandlers()

	err := c.g.Wait()
	c.stop()

	return err
}

func (c *impl) createAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	createIfNotExists := func(id string) {
		if _, exists := c.channels[id]; !exists {
			c.channels[id] = make(chan string, c.channelSize)
		}
	}

	for _, d := range c.decorators {
		createIfNotExists(d.input)
		createIfNotExists(d.output)
	}

	for _, m := range c.multiplexers {
		for _, input := range m.inputs {
			createIfNotExists(input)
		}
		createIfNotExists(m.output)
	}

	for _, s := range c.separators {
		createIfNotExists(s.input)
		for _, output := range s.outputs {
			createIfNotExists(output)
		}
	}
}

func (c *impl) startHandlers() {
	for _, d := range c.decorators {
		d := d
		c.g.Go(func() error {
			inputCh := c.getChannel(d.input)
			outputCh := c.getChannel(d.output)
			return d.fn(c.ctx, inputCh, outputCh)
		})
	}

	for _, m := range c.multiplexers {
		m := m
		c.g.Go(func() error {
			inputChs := make([]chan string, len(m.inputs))
			for i, inputID := range m.inputs {
				inputChs[i] = c.getChannel(inputID)
			}
			outputCh := c.getChannel(m.output)
			return m.fn(c.ctx, inputChs, outputCh)
		})
	}

	for _, s := range c.separators {
		s := s
		c.g.Go(func() error {
			inputCh := c.getChannel(s.input)
			outputChs := make([]chan string, len(s.outputs))
			for i, outputID := range s.outputs {
				outputChs[i] = c.getChannel(outputID)
			}
			return s.fn(c.ctx, inputCh, outputChs)
		})
	}
}

func (c *impl) stop() {
	c.mu.Lock()
	if c.stopped {
		c.mu.Unlock()
		return
	}
	c.stopped = true
	c.running = false
	c.mu.Unlock()

	if c.cancel != nil {
		c.cancel()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for id, ch := range c.channels {
		close(ch)
		delete(c.channels, id)
	}
}