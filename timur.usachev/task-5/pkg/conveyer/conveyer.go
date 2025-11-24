package conveyer

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

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
		decorators:   make([]decorator, 0),
		multiplexers: make([]multiplexer, 0),
		separators:   make([]separator, 0),
	}
}

func (c *Conveyer) RegisterDecorator(fn func(ctx context.Context, input, output chan string) error, inputName, outputName string) {
	if _, exists := c.channels[inputName]; !exists {
		c.channels[inputName] = make(chan string, c.size)
	}
	if _, exists := c.channels[outputName]; !exists {
		c.channels[outputName] = make(chan string, c.size)
	}
	c.decorators = append(c.decorators, decorator{
		fn:         fn,
		inputChan:  c.channels[inputName],
		outputChan: c.channels[outputName],
	})
}

func (c *Conveyer) RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputNames []string, outputName string) {
	inputChans := make([]chan string, len(inputNames))
	for i, name := range inputNames {
		if _, exists := c.channels[name]; !exists {
			c.channels[name] = make(chan string, c.size)
		}
		inputChans[i] = c.channels[name]
	}
	if _, exists := c.channels[outputName]; !exists {
		c.channels[outputName] = make(chan string, c.size)
	}
	c.multiplexers = append(c.multiplexers, multiplexer{
		fn:         fn,
		inputChans: inputChans,
		outputChan: c.channels[outputName],
	})
}

func (c *Conveyer) RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, inputName string, outputNames []string) {
	if _, exists := c.channels[inputName]; !exists {
		c.channels[inputName] = make(chan string, c.size)
	}
	inputChan := c.channels[inputName]
	outputChans := make([]chan string, len(outputNames))
	for i, name := range outputNames {
		if _, exists := c.channels[name]; !exists {
			c.channels[name] = make(chan string, c.size)
		}
		outputChans[i] = c.channels[name]
	}
	c.separators = append(c.separators, separator{
		fn:          fn,
		inputChan:   inputChan,
		outputChans: outputChans,
	})
}

func (c *Conveyer) Send(inputName string, data string) error {
	ch, exists := c.channels[inputName]
	if !exists {
		return errors.New("chan not found")
	}
	ch <- data
	return nil
}

func (c *Conveyer) Recv(outputName string) (string, error) {
	ch, exists := c.channels[outputName]
	if !exists {
		return "", errors.New("chan not found")
	}
	val, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return val, nil
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
		close(ch)
	}
	return err
}
