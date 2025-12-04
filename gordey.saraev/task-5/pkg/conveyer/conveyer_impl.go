package conveyer

import (
	"context"
	"errors"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound = errors.New("chan not found")
)

const undefined = "undefined"

type Conveyer struct {
	channels map[string]chan string
	handlers []func(ctx context.Context) error
	size     int
	mu       sync.Mutex
	running  bool
}

func New(size int) *Conveyer {
	if size < 0 {
		size = 0
	}
	return &Conveyer{
		channels: make(map[string]chan string),
		size:     size,
	}
}

func (c *Conveyer) getChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.channels[name]; ok {
		return ch
	}
	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(ctx context.Context, in chan string, out chan string) error,
	in string,
	out string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inCh := c.getChannel(in)
	outCh := c.getChannel(out)
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, ins []chan string, out chan string) error,
	ins []string,
	out string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inChs := make([]chan string, len(ins))
	for i, name := range ins {
		inChs[i] = c.getChannel(name)
	}
	outCh := c.getChannel(out)
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inChs, outCh)
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, in chan string, outs []chan string) error,
	in string,
	outs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inCh := c.getChannel(in)
	outChs := make([]chan string, len(outs))
	for i, name := range outs {
		outChs[i] = c.getChannel(name)
	}
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inCh, outChs)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return errors.New("already running")
	}
	c.running = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.running = false
		// Закрываем все каналы после завершения
		for name, ch := range c.channels {
			close(ch)
			delete(c.channels, name)
		}
		c.mu.Unlock()
	}()

	g, ctx := errgroup.WithContext(ctx)
	for _, h := range c.handlers {
		h := h
		g.Go(func() error {
			return h(ctx)
		})
	}
	return g.Wait()
}

func (c *Conveyer) Send(in string, data string) error {
	c.mu.Lock()
	ch, ok := c.channels[in]
	running := c.running
	c.mu.Unlock()

	if !running {
		return errors.New("conveyer not running")
	}
	if !ok {
		return ErrChanNotFound
	}

	select {
	case ch <- data:
		return nil
	default:
		return errors.New("channel full")
	}
}

func (c *Conveyer) Recv(out string) (string, error) {
	c.mu.Lock()
	ch, ok := c.channels[out]
	running := c.running
	c.mu.Unlock()

	if !running {
		return "", errors.New("conveyer not running")
	}
	if !ok {
		return "", ErrChanNotFound
	}

	val, ok := <-ch
	if !ok {
		return undefined, nil
	}
	return val, nil
}

// Функции обработчиков с исправлениями

var (
	ErrCantBeDecorated = errors.New("can't be decorated")
)

func PrefixDecoratorFunc(ctx context.Context, in chan string, out chan string) error {
	const prefix = "decorated: "
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-in:
			if !ok {
				close(out)
				return nil
			}
			if strings.Contains(data, "no decorator") {
				return ErrCantBeDecorated
			}
			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- data:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, in chan string, outs []chan string) error {
	if len(outs) == 0 {
		return nil
	}

	defer func() {
		for _, out := range outs {
			close(out)
		}
	}()

	i := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-in:
			if !ok {
				return nil
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case outs[i%len(outs)] <- data:
				i++
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, ins []chan string, out chan string) error {
	if len(ins) == 0 {
		close(out)
		return nil
	}

	defer close(out)

	done := make(chan struct{})
	defer close(done)

	type item struct {
		data string
		ok   bool
	}
	merged := make(chan item)

	for _, ch := range ins {
		ch := ch
		go func() {
			for {
				select {
				case <-done:
					return
				case <-ctx.Done():
					return
				case data, ok := <-ch:
					select {
					case <-done:
						return
					case <-ctx.Done():
						return
					case merged <- item{data, ok}:
						if !ok {
							return
						}
					}
				}
			}
		}()
	}

	alive := len(ins)
	for alive > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case it := <-merged:
			if !it.ok {
				alive--
				continue
			}
			if strings.Contains(it.data, "no multiplexer") {
				continue
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- it.data:
			}
		}
	}
	return nil
}
