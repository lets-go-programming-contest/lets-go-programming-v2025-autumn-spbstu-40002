package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	errSendChanNotFound = errors.New("chan not found")
	errRecvChanNotFound = errors.New("chan not found")
	errNoHandlers       = errors.New("conveyer has no handlers")
	undefined           = "undefined"
)

type Conveyer struct {
	channelSize  int
	channels     map[string]chan string
	handlersPool []func(context.Context) error
	mu           sync.RWMutex
}

func New(channelSize int) *Conveyer {
	return &Conveyer{
		channelSize:  channelSize,
		channels:     make(map[string]chan string),
		handlersPool: make([]func(context.Context) error, 0),
	}
}

func (c *Conveyer) makeChannel(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.channels[name]; !ok {
		c.channels[name] = make(chan string, c.channelSize)
	}
}

func (c *Conveyer) addToPool(function func(context.Context) error) {
	c.handlersPool = append(c.handlersPool, function)
}

func (c *Conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for name, ch := range c.channels {
		close(ch)
		delete(c.channels, name)
	}
}

func (c *Conveyer) RegisterDecorator(
	decoratorFunc func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.makeChannel(input)
	c.makeChannel(output)

	c.addToPool(func(context context.Context) error {
		c.mu.RLock()
		defer c.mu.RUnlock()

		inChan := c.channels[input]
		outChan := c.channels[output]

		return decoratorFunc(context, inChan, outChan)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	multiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	for _, input := range inputs {
		c.makeChannel(input)
	}
	c.makeChannel(output)

	c.addToPool(func(context context.Context) error {
		c.mu.RLock()
		defer c.mu.RUnlock()

		inputChannels := make([]chan string, len(inputs))
		outChan := c.channels[output]

		for i, inputName := range inputs {
			inputChannels[i] = c.channels[inputName]
		}

		return multiplexerFunc(
			context,
			inputChannels,
			outChan,
		)
	})
}

func (c *Conveyer) RegisterSeparator(
	separatorFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.makeChannel(input)
	for _, output := range outputs {
		c.makeChannel(output)
	}

	c.addToPool(func(context context.Context) error {
		c.mu.RLock()
		defer c.mu.RUnlock()

		inChan := c.channels[input]
		outputChannels := make([]chan string, len(outputs))

		for i, outputName := range outputs {
			outputChannels[i] = c.channels[outputName]
		}

		return separatorFunc(
			context,
			inChan,
			outputChannels,
		)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	if len(c.handlersPool) == 0 {
		return errNoHandlers
	}

	errChan := make(chan error, len(c.handlersPool))
	var wg sync.WaitGroup

	for _, handler := range c.handlersPool {
		wg.Add(1)
		go func(h func(context.Context) error) {
			defer wg.Done()
			if err := h(ctx); err != nil {
				errChan <- err
			}
		}(handler)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		c.closeAllChannels()

		return ctx.Err()
	case err, ok := <-errChan:
		if ok {
			c.closeAllChannels()
			<-done

			return err
		}
		c.closeAllChannels()
		<-done

		return nil
	}
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, exists := c.channels[input]
	if !exists {
		return errSendChanNotFound
	}

	ch <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, exists := c.channels[output]
	if !exists {
		return "", errRecvChanNotFound
	}

	data, ok := <-ch
	if !ok {
		return undefined, nil
	}

	return data, nil
}
