package conveyer

import (
    "context"
    "errors"
    "fmt"
    "sync"
    "golang.org/x/sync/errgroup"
)

var ErrChannelUnavailable = errors.New("channel unavailable")

const NoValuePlaceholder = "undefined"

type ConveyerInterface interface {
    RegisterDecorator(handler func(context.Context, chan string, chan string) error, inputName, outputName string)
    RegisterCombiner(handler func(context.Context, []chan string, chan string) error, inputNames []string, outputName string)
    RegisterSplitter(handler func(context.Context, chan string, []chan string) error, inputName string, outputNames []string)
    Run(ctx context.Context) error
    Send(input string, data string) error
    Recv(output string) (string, error)
}

type taskHandler func(context.Context) error

type Conveyer struct {
    bufferSize int
    channels   map[string]chan string
    tasks      []taskHandler
    mutex      sync.RWMutex
}

func New(size int) *Conveyer {
    return &Conveyer{
        bufferSize: size,
        channels:   make(map[string]chan string),
        tasks:      make([]taskHandler, 0),
        mutex:      sync.RWMutex{},
    }
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    if existingChan, found := c.channels[name]; found {
        return existingChan
    }

    newChan := make(chan string, c.bufferSize)
    c.channels[name] = newChan
    return newChan
}

func (c *Conveyer) getChannel(name string) (chan string, error) {
    c.mutex.RLock()
    defer c.mutex.RUnlock()

    if ch, found := c.channels[name]; found {
        return ch, nil
    }

    return nil, ErrChannelUnavailable
}

func (c *Conveyer) RegisterDecorator(
    handler func(context.Context, chan string, chan string) error,
    input, output string,
) {
    inChan := c.getOrCreateChannel(input)
    outChan := c.getOrCreateChannel(output)

    task := func(ctx context.Context) error {
        return handler(ctx, inChan, outChan)
    }

    c.addTask(task)
}

func (c *Conveyer) RegisterCombiner(
    handler func(context.Context, []chan string, chan string) error,
    inputs []string,
    output string,
) {
    inChannels := make([]chan string, len(inputs))
    for i, name := range inputs {
        inChannels[i] = c.getOrCreateChannel(name)
    }

    outChan := c.getOrCreateChannel(output)

    task := func(ctx context.Context) error {
        return handler(ctx, inChannels, outChan)
    }

    c.addTask(task)
}

func (c *Conveyer) RegisterSplitter(
    handler func(context.Context, chan string, []chan string) error,
    input string,
    outputs []string,
) {
    inChan := c.getOrCreateChannel(input)

    outChannels := make([]chan string, len(outputs))
    for i, name := range outputs {
        outChannels[i] = c.getOrCreateChannel(name)
    }

    task := func(ctx context.Context) error {
        return handler(ctx, inChan, outChannels)
    }

    c.addTask(task)
}

func (c *Conveyer) Send(input string, data string) error {
    ch, err := c.getChannel(input)
    if err != nil {
        return err
    }

    select {
    case ch <- data:
        return nil
    default:
        return fmt.Errorf("channel buffer is full for %s", input)
    }
}

func (c *Conveyer) Recv(output string) (string, error) {
    ch, err := c.getChannel(output)
    if err != nil {
        return "", err
    }

    data, ok := <-ch
    if !ok {
        return NoValuePlaceholder, nil
    }

    return data, nil
}

func (c *Conveyer) closeChannels() {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    for _, ch := range c.channels {
        close(ch)
    }
}

func (c *Conveyer) Run(ctx context.Context) error {
    c.mutex.RLock()
    defer c.mutex.RUnlock()

    group, ctx := errgroup.WithContext(ctx)

    for _, task := range c.tasks {
        taskCopy := task
        group.Go(func() error {
            return taskCopy(ctx)
        })
    }

    err := group.Wait()
    c.closeChannels()

    if err != nil {
        return fmt.Errorf("conveyer pipeline failed: %w", err)
    }

    return nil
}

func (c *Conveyer) addTask(task taskHandler) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.tasks = append(c.tasks, task)
}
