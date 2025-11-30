package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

const undefined = "undefined"

type handlerFunc func(ctx context.Context) error

type Conveyer struct {
	size     int
	mu       sync.RWMutex
	chans    map[string]chan string
	handlers []handlerFunc
}

// New creates a new conveyer instance with specified channel buffer size.
func New(size int) *Conveyer {
	if size < 0 {
		size = 0
	}

	return &Conveyer{
		size:     size,
		mu:       sync.RWMutex{},
		chans:    make(map[string]chan string),
		handlers: nil,
	}
}

// ensureChan creates a channel if it doesn't exist.
func (c *Conveyer) ensureChan(chanID string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if existingChan, exists := c.chans[chanID]; exists {
		return existingChan
	}

	createdChan := make(chan string, c.size)
	c.chans[chanID] = createdChan

	return createdChan
}

// getChan retrieves a channel by ID.
func (c *Conveyer) getChan(chanID string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, exists := c.chans[chanID]

	return ch, exists
}

// RegisterDecorator registers a data modifier handler.
func (c *Conveyer) RegisterDecorator(
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

// RegisterMultiplexer registers a multiplexer handler.
func (c *Conveyer) RegisterMultiplexer(
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

// RegisterSeparator registers a separator handler.
func (c *Conveyer) RegisterSeparator(
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

// Run starts the conveyer and runs all handlers in separate goroutines.
func (c *Conveyer) Run(ctx context.Context) error {
	if len(c.handlers) == 0 {
		return nil
	}

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
				!errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
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

// Send sends data to a channel identified by input ID.
func (c *Conveyer) Send(input string, data string) error {
	channelInstance, exists := c.getChan(input)
	if !exists {
		return fmt.Errorf("send failed: %w", ErrChannelNotFound)
	}

	channelInstance <- data

	return nil
}

// Recv receives data from a channel identified by output ID.
func (c *Conveyer) Recv(output string) (string, error) {
	channelInstance, exists := c.getChan(output)
	if !exists {
		return "", fmt.Errorf("recv failed: %w", ErrChannelNotFound)
	}

	value, okChannel := <-channelInstance
	if !okChannel {
		return undefined, nil
	}

	return value, nil
}
