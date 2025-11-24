package conveyer

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

type decorator struct {
	fn         func(ctx context.Context, input, output chan string) error
	inputChan  chan string
	outputChan chan string
}

type multiplexer struct {
	fn         func(ctx context.Context, inputs []chan string, output chan string) error
	inputChans []chan string
	outputChan chan string
}

type separator struct {
	fn          func(ctx context.Context, input chan string, outputs []chan string) error
	inputChan   chan string
	outputChans []chan string
}

type Conveyer struct {
	size         int
	channels     map[string]chan string
	decorators   []decorator
	multiplexers []multiplexer
	separators   []separator
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:         size,
		channels:     make(map[string]chan string),
		decorators:   nil,
		multiplexers: nil,
		separators:   nil,
	}
}

func (c *Conveyer) RegisterDecorator(fn func(ctx context.Context, input, output chan string) error, inputName, outputName string) {
	if _, ok := c.channels[inputName]; !ok {
		c.channels[inputName] = make(chan string, c.size)
	}
	if _, ok := c.channels[outputName]; !ok {
		c.channels[outputName] = make(chan string, c.size)
	}
	c.decorators = append(c.decorators, decorator{
		fn:         fn,
		inputChan:  c.channels[inputName],
		outputChan: c.channels[outputName],
	})
}

func (c *Conveyer) RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputNames []string, outputName string) {
	inputs := make([]chan string, len(inputNames))
	for i, n := range inputNames {
		if _, ok := c.channels[n]; !ok {
			c.channels[n] = make(chan string, c.size)
		}
		inputs[i] = c.channels[n]
	}
	if _, ok := c.channels[outputName]; !ok {
		c.channels[outputName] = make(chan string, c.size)
	}
	c.multiplexers = append(c.multiplexers, multiplexer{
		fn:         fn,
		inputChans: inputs,
		outputChan: c.channels[outputName],
	})
}

func (c *Conveyer) RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, inputName string, outputNames []string) {
	if _, ok := c.channels[inputName]; !ok {
		c.channels[inputName] = make(chan string, c.size)
	}
	input := c.channels[inputName]
	outChans := make([]chan string, len(outputNames))
	for i, n := range outputNames {
		if _, ok := c.channels[n]; !ok {
			c.channels[n] = make(chan string, c.size)
		}
		outChans[i] = c.channels[n]
	}
	c.separators = append(c.separators, separator{
		fn:          fn,
		inputChan:   input,
		outputChans: outChans,
	})
}

func (c *Conveyer) Send(inputName string, data string) error {
	ch, ok := c.channels[inputName]
	if !ok {
		return ErrChanNotFound
	}
	ch <- data
	return nil
}

func (c *Conveyer) Recv(outputName string) (string, error) {
	ch, ok := c.channels[outputName]
	if !ok {
		return "", ErrChanNotFound
	}
	v, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return v, nil
}

func safeClose(ch chan string) {
	if ch == nil {
		return
	}
	defer func() { _ = recover() }()
	close(ch)
}

func (c *Conveyer) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, d := range c.decorators {
		d := d
		g.Go(func() error {
			return d.fn(ctx, d.inputChan, d.outputChan)
		})
	}

	for _, m := range c.multiplexers {
		m := m
		g.Go(func() error {
			return m.fn(ctx, m.inputChans, m.outputChan)
		})
	}

	for _, s := range c.separators {
		s := s
		g.Go(func() error {
			return s.fn(ctx, s.inputChan, s.outputChans)
		})
	}

	err := g.Wait()

	for _, ch := range c.channels {
		safeClose(ch)
	}
	return err
}
