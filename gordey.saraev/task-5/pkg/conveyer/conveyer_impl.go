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
	chans map[string]chan string
	funcs []func(context.Context) error
}

func New(size int) *Conveyer {
	return &Conveyer{
		chans: make(map[string]chan string),
	}
}

func (c *Conveyer) getChan(name string) chan string {
	if ch, ok := c.chans[name]; ok {
		return ch
	}
	ch := make(chan string, 100)
	c.chans[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(fn func(context.Context, chan string, chan string) error, in, out string) {
	inCh := c.getChan(in)
	outCh := c.getChan(out)
	c.funcs = append(c.funcs, func(ctx context.Context) error { return fn(ctx, inCh, outCh) })
}

func (c *Conveyer) RegisterMultiplexer(fn func(context.Context, []chan string, chan string) error, ins []string, out string) {
	inChs := make([]chan string, len(ins))
	for i, name := range ins {
		inChs[i] = c.getChan(name)
	}
	outCh := c.getChan(out)
	c.funcs = append(c.funcs, func(ctx context.Context) error { return fn(ctx, inChs, outCh) })
}

func (c *Conveyer) RegisterSeparator(fn func(context.Context, chan string, []chan string) error, in string, outs []string) {
	inCh := c.getChan(in)
	outChs := make([]chan string, len(outs))
	for i, name := range outs {
		outChs[i] = c.getChan(name)
	}
	c.funcs = append(c.funcs, func(ctx context.Context) error { return fn(ctx, inCh, outChs) })
}

func (c *Conveyer) Run(ctx context.Context) error {
	gr, ctx := errgroup.WithContext(ctx)
	for _, f := range c.funcs {
		f := f
		gr.Go(func() error { return f(ctx) })
	}
	return gr.Wait()
}

func (c *Conveyer) Send(in string, data string) error {
	ch, ok := c.chans[in]
	if !ok {
		return ErrChanNotFound
	}
	ch <- data
	return nil
}

func (c *Conveyer) Recv(out string) (string, error) {
	ch, ok := c.chans[out]
	if !ok {
		return "", ErrChanNotFound
	}
	val, ok := <-ch
	if !ok {
		return undefined, nil
	}
	return val, nil
}

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
				return errors.New("can't be decorated")
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
