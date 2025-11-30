package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")
var ErrAlreadyClosed = errors.New("already closed")

type conveyer interface {
	RegisterDecorator(
		fn func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		fn func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		fn func(ctx context.Context, input chan string, outputs []chan string) error,
		input string,
		outputs []string,
	)

	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type decoratorReg struct {
	fn  func(ctx context.Context, input chan string, output chan string) error
	in  string
	out string
}

type multiplexerReg struct {
	fn  func(ctx context.Context, inputs []chan string, output chan string) error
	ins []string
	out string
}

type separatorReg struct {
	fn   func(ctx context.Context, input chan string, outputs []chan string) error
	in   string
	outs []string
}

type conveyerImpl struct {
	size         int
	mu           sync.RWMutex
	chans        map[string]chan string
	decorators   []decoratorReg
	multiplexers []multiplexerReg
	separators   []separatorReg

	closed bool
}

func New(size int) conveyer {
	return &conveyerImpl{
		size:  size,
		chans: make(map[string]chan string),
	}
}

func (c *conveyerImpl) ensureChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ch, ok := c.chans[name]; ok {
		return ch
	}
	ch := make(chan string, c.size)
	c.chans[name] = ch
	return ch
}

func (c *conveyerImpl) RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, input string, output string) {
	c.ensureChan(input)
	c.ensureChan(output)
	c.decorators = append(c.decorators, decoratorReg{fn: fn, in: input, out: output})
}

func (c *conveyerImpl) RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string) {
	for _, n := range inputs {
		c.ensureChan(n)
	}
	c.ensureChan(output)
	c.multiplexers = append(c.multiplexers, multiplexerReg{fn: fn, ins: inputs, out: output})
}

func (c *conveyerImpl) RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string) {
	c.ensureChan(input)
	for _, n := range outputs {
		c.ensureChan(n)
	}
	c.separators = append(c.separators, separatorReg{fn: fn, in: input, outs: outputs})
}

func (c *conveyerImpl) getChan(name string) chan string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.chans[name]
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return ErrAlreadyClosed
	}
	c.mu.RUnlock()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	launch := func(fn func(context.Context) error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fn(ctx); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				select {
				case errCh <- err:
				default:
				}
			}
		}()
	}

	for _, d := range c.decorators {
		in := c.getChan(d.in)
		out := c.getChan(d.out)
		fn := d.fn
		launch(func(cctx context.Context) error {
			select {
			case <-cctx.Done():
				return cctx.Err()
			default:
				return fn(cctx, in, out)
			}
		})
	}

	for _, m := range c.multiplexers {
		out := c.getChan(m.out)
		ins := make([]chan string, len(m.ins))
		for i, n := range m.ins {
			ins[i] = c.getChan(n)
		}
		fn := m.fn
		launch(func(cctx context.Context) error {
			select {
			case <-cctx.Done():
				return cctx.Err()
			default:
				return fn(cctx, ins, out)
			}
		})
	}

	for _, s := range c.separators {
		in := c.getChan(s.in)
		outs := make([]chan string, len(s.outs))
		for i, n := range s.outs {
			outs[i] = c.getChan(n)
		}
		fn := s.fn
		launch(func(cctx context.Context) error {
			select {
			case <-cctx.Done():
				return cctx.Err()
			default:
				return fn(cctx, in, outs)
			}
		})
	}

	wg.Wait()

	c.mu.Lock()
	if !c.closed {
		for _, ch := range c.chans {
			close(ch)
		}
		c.closed = true
	}
	c.mu.Unlock()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func (c *conveyerImpl) closeAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}

	for _, ch := range c.chans {
		func() {
			defer func() { recover() }()
			close(ch)
		}()
	}

	c.closed = true
}

func (c *conveyerImpl) Send(input string, data string) (err error) {
	c.mu.RLock()
	ch, ok := c.chans[input]
	c.mu.RUnlock()
	if !ok {
		return ErrChanNotFound
	}
	defer func() {
		if r := recover(); r != nil {
			err = ErrChanNotFound
		}
	}()
	ch <- data
	return nil
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, ok := c.chans[output]
	c.mu.RUnlock()
	if !ok {
		return "", ErrChanNotFound
	}
	v, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return v, nil
}
