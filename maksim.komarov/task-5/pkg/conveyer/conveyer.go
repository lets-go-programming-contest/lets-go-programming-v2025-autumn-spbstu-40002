package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

const chanNotFoundMsg = "chan not found"

var ErrChannelNotFound = errors.New(chanNotFoundMsg)

type handlerFunc func(ctx context.Context) error

type stringConveyer struct {
	size int

	mu       sync.RWMutex
	chans    map[string]chan string
	handlers []handlerFunc
}

func New(size int) *stringConveyer {
	if size < 0 {
		size = 0
	}

	return &stringConveyer{
		size:     size,
		mu:       sync.RWMutex{},
		chans:    make(map[string]chan string),
		handlers: nil,
	}
}

func (c *stringConveyer) ensureChan(chanID string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.chans[chanID]; exists {
		return ch
	}

	createdChan := make(chan string, c.size)
	c.chans[chanID] = createdChan

	return createdChan
}

func (c *stringConveyer) getChan(chanID string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, exists := c.chans[chanID]

	return ch, exists
}

func (c *stringConveyer) RegisterDecorator(
	handler func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputChan := c.ensureChan(input)
	outputChan := c.ensureChan(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, inputChan, outputChan)
	})
}

func (c *stringConveyer) RegisterMultiplexer(
	handler func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputChans := make([]chan string, 0, len(inputs))

	for _, chanID := range inputs {
		inputChans = append(inputChans, c.ensureChan(chanID))
	}

	outputChan := c.ensureChan(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, inputChans, outputChan)
	})
}

func (c *stringConveyer) RegisterSeparator(
	handler func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputChan := c.ensureChan(input)
	outputChans := make([]chan string, 0, len(outputs))

	for _, chanID := range outputs {
		outputChans = append(outputChans, c.ensureChan(chanID))
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, inputChan, outputChans)
	})
}

func (c *stringConveyer) Run(ctx context.Context) error {
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var waitGroup sync.WaitGroup

	errorChan := make(chan error, 1)

	for _, registeredHandler := range c.handlers {
		handlerCopy := registeredHandler

		waitGroup.Add(1)

		runHandler := func() {
			defer waitGroup.Done()

			if err := handlerCopy(runCtx); err != nil &&
				!errors.Is(err, context.Canceled) &&
				!errors.Is(err, context.DeadlineExceeded) {
				select {
				case errorChan <- err:
				default:
				}
			}
		}

		go runHandler()
	}

	var runError error

	select {
	case <-ctx.Done():
		runError = nil
	case err := <-errorChan:
		runError = err
	}

	cancel()
	waitGroup.Wait()

	c.mu.Lock()

	for _, channelInstance := range c.chans {
		close(channelInstance)
	}

	c.mu.Unlock()

	return runError
}

func (c *stringConveyer) Send(input string, data string) error {
	channelInstance, exists := c.getChan(input)
	if !exists {
		return fmt.Errorf("send failed: %w", ErrChannelNotFound)
	}

	channelInstance <- data

	return nil
}

func (c *stringConveyer) Recv(output string) (string, error) {
	channelInstance, exists := c.getChan(output)
	if !exists {
		return "", fmt.Errorf("recv failed: %w", ErrChannelNotFound)
	}

	value, okChannel := <-channelInstance
	if !okChannel {
		return "undefined", nil
	}

	return value, nil
}
