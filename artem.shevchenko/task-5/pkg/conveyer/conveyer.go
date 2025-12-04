package conveyer

import (
	"context"
	"errors"
)

const undefinedData = "undefined"

var ErrChanNotFound = errors.New("chan not found")

type impl struct {
	size     int
	channels map[string]chan string
	tasks    []func(context.Context) error
}

func New(size int) conveyer {
	return &impl{
		size:     size,
		channels: make(map[string]chan string),
		tasks:    make([]func(context.Context) error, 0),
	}
}

func (c *impl) getOrCreateChan(name string) chan string {
	if ch, ok := c.channels[name]; ok {
		return ch
	}
	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

// === Методы интерфейса ===

func (c *impl) RegisterDecorator(
	fn func(ctx context.Context, input, output chan string) error,
	input, output string,
) {
	in := c.getOrCreateChan(input)
	out := c.getOrCreateChan(output)
	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *impl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChans[i] = c.getOrCreateChan(name)
	}
	out := c.getOrCreateChan(output)
	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return fn(ctx, inChans, out)
	})
}

func (c *impl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	in := c.getOrCreateChan(input)
	outChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChans[i] = c.getOrCreateChan(name)
	}
	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return fn(ctx, in, outChans)
	})
}

func (c *impl) Run(ctx context.Context) error {
	if len(c.tasks) == 0 {
		<-ctx.Done()
		return ctx.Err()
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 1)

	for _, task := range c.tasks {
		task := task
		go func() {
			if err := task(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
				cancel()
			}
		}()
	}

	select {
	case err := <-errCh:
		for _, ch := range c.channels {
			close(ch)
		}
		return err
	case <-ctx.Done():
		for _, ch := range c.channels {
			close(ch)
		}
		return ctx.Err()
	}
}

func (c *impl) Send(name, data string) error {
	ch, ok := c.channels[name]
	if !ok {
		return ErrChanNotFound
	}
	ch <- data
	return nil
}

func (c *impl) Recv(name string) (string, error) {
	ch, ok := c.channels[name]
	if !ok {
		return "", ErrChanNotFound
	}
	data, ok := <-ch
	if !ok {
		return undefinedData, nil // по ТЗ — НЕ ошибка!
	}
	return data, nil
}