package conveyer

import (
	"context"
	"errors"
	"sync"
)

type DecoratorFunc func(
	ctx context.Context,
	input chan string,
	output chan string,
) error

type MultiplexerFunc func(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error

type SeparatorFunc func(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error

type Conveyer interface {
	RegisterDecorator(fn DecoratorFunc, input string, output string) error
	RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string) error
	RegisterSeparator(fn SeparatorFunc, input string, outputs []string) error

	Run(ctx context.Context) error

	Send(input string, data string) error
	Recv(output string) (string, error)
}

const Undefined = "undefined"

var ErrChanNotFound = errors.New("chan not found")

type runFn func(context.Context) error

type conv struct {
	size int

	mu      sync.RWMutex
	chans   map[string]chan string
	runners []runFn

	startMux sync.Mutex
	started  bool
}

func New(size int) Conveyer {
	return &conv{
		size:    size,
		chans:   make(map[string]chan string),
		runners: make([]runFn, 0),
	}
}

func (c *conv) ensureChan(id string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	exists, ok := c.chans[id]
	if ok {
		return exists
	}

	created := make(chan string, c.size)
	c.chans[id] = created

	return created
}

func (c *conv) getChan(id string) (chan string, bool) {
	c.mu.RLock()
	ch, ok := c.chans[id]
	c.mu.RUnlock()

	return ch, ok
}

func (c *conv) RegisterDecorator(
	fn DecoratorFunc,
	input string,
	output string,
) error {
	inputChan := c.ensureChan(input)
	outputChan := c.ensureChan(output)

	c.runners = append(c.runners, func(ctx context.Context) error {
		return fn(ctx, inputChan, outputChan)
	})

	return nil
}

func (c *conv) RegisterMultiplexer(
	fn MultiplexerFunc,
	inputs []string,
	output string,
) error {
	inputChans := make([]chan string, 0, len(inputs))

	for _, id := range inputs {
		inputChans = append(inputChans, c.ensureChan(id))
	}

	outputChan := c.ensureChan(output)

	c.runners = append(c.runners, func(ctx context.Context) error {
		return fn(ctx, inputChans, outputChan)
	})

	return nil
}

func (c *conv) RegisterSeparator(
	fn SeparatorFunc,
	input string,
	outputs []string,
) error {
	inputChan := c.ensureChan(input)

	outputChans := make([]chan string, 0, len(outputs))

	for _, id := range outputs {
		outputChans = append(outputChans, c.ensureChan(id))
	}

	c.runners = append(c.runners, func(ctx context.Context) error {
		return fn(ctx, inputChan, outputChans)
	})

	return nil
}

func (c *conv) closeAllChans() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, ch := range c.chans {
		func() {
			defer func() { _ = recover() }()
			close(ch)
		}()
		delete(c.chans, id)
	}
}

func (c *conv) Run(ctx context.Context) error {
	c.startMux.Lock()
	already := c.started
	if !already {
		c.started = true
	}
	c.startMux.Unlock()

	if already {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(c.runners))

	for _, rf := range c.runners {
		wg.Add(1)

		go func(f runFn) {
			defer wg.Done()

			if err := f(ctx); err != nil {
				errChan <- err
			}
		}(rf)
	}

	waitDone := make(chan struct{})

	go func() {
		wg.Wait()
		close(waitDone)
	}()

	var firstErr error

	select {
	case firstErr = <-errChan:

	case <-waitDone:
	}

	c.closeAllChans()

	return firstErr
}

func (c *conv) Send(input string, data string) error {
	ch, ok := c.getChan(input)
	if !ok {
		return ErrChanNotFound
	}

	ch <- data

	return nil
}

func (c *conv) Recv(output string) (string, error) {
	ch, ok := c.getChan(output)
	if !ok {
		return "", ErrChanNotFound
	}

	value, ok := <-ch
	if !ok {
		return Undefined, nil
	}

	return value, nil
}
