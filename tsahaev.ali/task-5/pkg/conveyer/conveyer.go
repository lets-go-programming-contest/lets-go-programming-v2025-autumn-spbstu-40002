package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrChanNotFound   = errors.New("chan not found")
	ErrChanFull       = errors.New("chan is full")
	ErrNotInitialized = errors.New("conveyer not initialized")
)

const StrUndefined = "undefined"

type Conveyer struct {
	channels     map[string]chan string
	size         int
	decorators   []decorator
	multiplexers []multiplexer
	separators   []separator
	mu           sync.Mutex
	initialized  bool
}

type decorator struct {
	function func(ctx context.Context, input chan string, output chan string) error
	in       string
	out      string
}

type multiplexer struct {
	function func(ctx context.Context, inputs []chan string, output chan string) error
	ins      []string
	out      string
}

type separator struct {
	function func(ctx context.Context, input chan string, outputs []chan string) error
	in       string
	outs     []string
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels:     make(map[string]chan string),
		size:         size,
		decorators:   make([]decorator, 0),
		multiplexers: make([]multiplexer, 0),
		separators:   make([]separator, 0),
		mu:           sync.Mutex{},
		initialized:  true,
	}
}

func (c *Conveyer) getOrCreateChan(channelID string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, exists := c.channels[channelID]; exists {
		return channel
	}

	channel := make(chan string, c.size)
	c.channels[channelID] = channel

	return channel
}

func (c *Conveyer) getChan(channelID string) (chan string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	channel, exists := c.channels[channelID]

	return channel, exists
}

func (c *Conveyer) RegisterDecorator(
	function func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.decorators = append(c.decorators, decorator{
		function: function,
		in:       input,
		out:      output,
	})

	c.getOrCreateChan(input)
	c.getOrCreateChan(output)
}

func (c *Conveyer) RegisterMultiplexer(
	function func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.multiplexers = append(c.multiplexers, multiplexer{
		function: function,
		ins:      inputs,
		out:      output,
	})

	for _, input := range inputs {
		c.getOrCreateChan(input)
	}

	c.getOrCreateChan(output)
}

func (c *Conveyer) RegisterSeparator(
	function func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.separators = append(c.separators, separator{
		function: function,
		in:       input,
		outs:     outputs,
	})

	c.getOrCreateChan(input)

	for _, output := range outputs {
		c.getOrCreateChan(output)
	}
}

func (c *Conveyer) Run(parentCtx context.Context) error {
	if !c.initialized {
		return ErrNotInitialized
	}

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	waitGroup := &sync.WaitGroup{}
	errorChannel := make(chan error, 1)

	c.startDecorators(ctx, waitGroup, errorChannel)
	c.startMultiplexers(ctx, waitGroup, errorChannel)
	c.startSeparators(ctx, waitGroup, errorChannel)
	c.startCompletionWatcher(waitGroup, errorChannel)

	return c.waitForCompletion(ctx, cancel, errorChannel)
}

func (c *Conveyer) startDecorators(ctx context.Context, waitGroup *sync.WaitGroup, errorChannel chan error) {
	for _, decoratorItem := range c.decorators {
		waitGroup.Add(1)

		go func(dec decorator) {
			defer waitGroup.Done()

			inputChannel := c.getOrCreateChan(dec.in)
			outputChannel := c.getOrCreateChan(dec.out)

			if err := dec.function(ctx, inputChannel, outputChannel); err != nil {
				c.sendError(errorChannel, err)
			}
		}(decoratorItem)
	}
}

func (c *Conveyer) startMultiplexers(ctx context.Context, waitGroup *sync.WaitGroup, errorChannel chan error) {
	for _, multiplexerItem := range c.multiplexers {
		waitGroup.Add(1)

		go func(mux multiplexer) {
			defer waitGroup.Done()

			inputChannels := make([]chan string, len(mux.ins))
			for index, input := range mux.ins {
				inputChannels[index] = c.getOrCreateChan(input)
			}

			outputChannel := c.getOrCreateChan(mux.out)

			if err := mux.function(ctx, inputChannels, outputChannel); err != nil {
				c.sendError(errorChannel, err)
			}
		}(multiplexerItem)
	}
}

func (c *Conveyer) startSeparators(ctx context.Context, waitGroup *sync.WaitGroup, errorChannel chan error) {
	for _, separatorItem := range c.separators {
		waitGroup.Add(1)

		go func(sep separator) {
			defer waitGroup.Done()

			inputChannel := c.getOrCreateChan(sep.in)
			outputChannels := make([]chan string, len(sep.outs))

			for index, output := range sep.outs {
				outputChannels[index] = c.getOrCreateChan(output)
			}

			if err := sep.function(ctx, inputChannel, outputChannels); err != nil {
				c.sendError(errorChannel, err)
			}
		}(separatorItem)
	}
}

func (c *Conveyer) startCompletionWatcher(waitGroup *sync.WaitGroup, errorChannel chan error) {
	go func() {
		waitGroup.Wait()
		c.closeAllChans()
		close(errorChannel)
	}()
}

func (c *Conveyer) waitForCompletion(ctx context.Context, cancel context.CancelFunc, errorChannel chan error) error {
	select {
	case err := <-errorChannel:
		cancel()

		return err
	case <-ctx.Done():
		return nil
	}
}

func (c *Conveyer) sendError(errorChannel chan error, err error) {
	select {
	case errorChannel <- err:
	default:
	}
}

func (c *Conveyer) closeAllChans() {
	c.mu.Lock()
	defer c.mu.Unlock()

	closed := make(map[chan string]bool)
	for _, channel := range c.channels {
		if !closed[channel] {
			close(channel)

			closed[channel] = true
		}
	}
}

func (c *Conveyer) Send(channelID string, data string) error {
	channel, exists := c.getChan(channelID)
	if !exists {
		return ErrChanNotFound
	}

	select {
	case channel <- data:
		return nil
	default:
		return ErrChanFull
	}
}

func (c *Conveyer) Recv(channelID string) (string, error) {
	channel, exists := c.getChan(channelID)
	if !exists {
		return "", ErrChanNotFound
	}

	data, ok := <-channel
	if !ok {
		return StrUndefined, nil
	}

	return data, nil
}
