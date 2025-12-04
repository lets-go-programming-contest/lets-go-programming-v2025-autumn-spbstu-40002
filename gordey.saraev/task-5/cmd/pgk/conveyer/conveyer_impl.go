package conveyer

import (
	"context"
	"errors"
	"strings"

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
	ch, ok := c.channels[in]
	if !ok {
		return ErrChanNotFound
	}
	ch <- data
	return nil
}

func (c *Conveyer) Recv(out string) (string, error) {
	ch, ok := c.channels[out]
	if !ok {
		return "", ErrChanNotFound
	}
	val, ok := <-ch
	if !ok {
		return undefined, nil
	}
	return val, nil
}

// Функции обработчиков (твои оригинальные)

var (
	ErrCantBeDecorated = errors.New("can't be decorated")
)

func PrefixDecoratorFunc(ctx context.Context, in chan string, out chan string) error {
	const prefix = "decorated: "
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-in:
			if !ok {
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
				return nil
			case out <- data:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, in chan string, outs []chan string) error {
	if len(outs) == 0 {
		return nil
	}
	i := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-in:
			if !ok {
				return nil
			}
			select {
			case <-ctx.Done():
				return nil
			case outs[i%len(outs)] <- data:
				i++
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, ins []chan string, out chan string) error {
	if len(ins) == 0 {
		return nil
	}
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
			return nil
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
				return nil
			case out <- it.data:
			}
		}
	}
	return nil
}
