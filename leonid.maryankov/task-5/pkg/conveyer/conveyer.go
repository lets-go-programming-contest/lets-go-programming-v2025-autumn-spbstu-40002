package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

type HandlerRunner func(ctx context.Context) error

type Channel struct {
	Ch chan string
}

type ConveyorImpl struct {
	Size    int
	Mu      sync.RWMutex
	Chans   map[string]*Channel
	Runners []HandlerRunner
}

var ErrChanNotFound = errors.New("chan not found")

func New(size int) *ConveyorImpl {
	return &ConveyorImpl{
		Size:    size,
		Mu:      sync.RWMutex{},
		Chans:   make(map[string]*Channel),
		Runners: make([]HandlerRunner, 0),
	}
}

func (c *ConveyorImpl) getOrCreate(channelID string) chan string {
	if channelID == "" {
		return nil
	}

	c.Mu.Lock()
	defer c.Mu.Unlock()

	existing, exists := c.Chans[channelID]
	if !exists {
		newChan := make(chan string, c.Size)
		c.Chans[channelID] = &Channel{Ch: newChan}

		return newChan
	}

	return existing.Ch
}

func (c *ConveyorImpl) RegisterDecorator(
	handlerFunc func(context.Context, chan string, chan string) error,
	inputID string,
	outputID string,
) {
	inputCh := c.getOrCreate(inputID)
	outputCh := c.getOrCreate(outputID)

	c.Runners = append(c.Runners, func(ctx context.Context) error {
		return handlerFunc(ctx, inputCh, outputCh)
	})
}

func (c *ConveyorImpl) RegisterSeparator(
	handlerFunc func(context.Context, chan string, []chan string) error,
	inputID string,
	outputIDs []string,
) {
	inputCh := c.getOrCreate(inputID)

	outputChs := make([]chan string, 0, len(outputIDs))
	for _, id := range outputIDs {
		outputChs = append(outputChs, c.getOrCreate(id))
	}

	c.Runners = append(c.Runners, func(ctx context.Context) error {
		return handlerFunc(ctx, inputCh, outputChs)
	})
}

func (c *ConveyorImpl) RegisterMultiplexer(
	handlerFunc func(context.Context, []chan string, chan string) error,
	inputIDs []string,
	outputID string,
) {
	inputChs := make([]chan string, 0, len(inputIDs))
	for _, id := range inputIDs {
		inputChs = append(inputChs, c.getOrCreate(id))
	}

	outputCh := c.getOrCreate(outputID)

	c.Runners = append(c.Runners, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChs, outputCh)
	})
}

func (c *ConveyorImpl) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for i := range c.Runners {
		runner := c.Runners[i]
		runnerLocal := runner

		group.Go(func() error {
			return runnerLocal(ctx)
		})
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("%w", err)
	}

	c.closeAll()

	return nil
}

func (c *ConveyorImpl) closeAll() {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	for _, chObj := range c.Chans {
		func(ch chan string) {
			defer func() {
				_ = recover()
			}()
			close(ch)
		}(chObj.Ch)
	}
}

func (c *ConveyorImpl) Send(inputID string, data string) error {
	c.Mu.RLock()
	channelObj, exists := c.Chans[inputID]
	c.Mu.RUnlock()

	if !exists {
		return ErrChanNotFound
	}

	channelObj.Ch <- data

	return nil
}

func (c *ConveyorImpl) Recv(outputID string) (string, error) {
	c.Mu.RLock()
	channelObj, exists := c.Chans[outputID]
	c.Mu.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	value, open := <-channelObj.Ch
	if !open {
		return "undefined", nil
	}

	return value, nil
}
