package conveyer

import (
    "context"
    "errors"
    "sync"
)

type Conveyer struct {
    mutex           sync.RWMutex
    channelSize     int
    channels        map[string]chan string
    handlersPool    []func(context.Context) error
    ctx             context.Context
    cancel          context.CancelFunc
    running         bool
}

func New(size int) *Conveyer {
    return &Conveyer{
        channelSize:  size,
        channels:     make(map[string]chan string),
        handlersPool: make([]func(context.Context) error, 0),
        running:      false,
    }
}

func (conveyer *Conveyer) getOrCreateChannel(name string) chan string {
    conveyer.mutex.Lock()
    defer conveyer.mutex.Unlock()
    
    if channel, exists := conveyer.channels[name]; exists {
        return channel
    }
    
    channel := make(chan string, conveyer.channelSize)
    conveyer.channels[name] = channel
    return channel
}

func (conveyer *Conveyer) RegisterDecorator(
    fn func(ctx context.Context, input chan string, output chan string) error,
    input string,
    output string,
) {
    inputChannel := conveyer.getOrCreateChannel(input)
    outputChannel := conveyer.getOrCreateChannel(output)
    
    conveyer.mutex.Lock()
    conveyer.handlersPool = append(conveyer.handlersPool, func(ctx context.Context) error {
        return fn(ctx, inputChannel, outputChannel)
    })
    conveyer.mutex.Unlock()
}

func (conveyer *Conveyer) RegisterMultiplexer(
    fn func(ctx context.Context, inputs []chan string, output chan string) error,
    inputs []string,
    output string,
) {
    inputChannels := make([]chan string, len(inputs))
    for index, input := range inputs {
        inputChannels[index] = conveyer.getOrCreateChannel(input)
    }
    outputChannel := conveyer.getOrCreateChannel(output)
    
    conveyer.mutex.Lock()
    conveyer.handlersPool = append(conveyer.handlersPool, func(ctx context.Context) error {
        return fn(ctx, inputChannels, outputChannel)
    })
    conveyer.mutex.Unlock()
}

func (conveyer *Conveyer) RegisterSeparator(
    fn func(ctx context.Context, input chan string, outputs []chan string) error,
    input string,
    outputs []string,
) {
    inputChannel := conveyer.getOrCreateChannel(input)
    outputChannels := make([]chan string, len(outputs))
    for index, output := range outputs {
        outputChannels[index] = conveyer.getOrCreateChannel(output)
    }
    
    conveyer.mutex.Lock()
    conveyer.handlersPool = append(conveyer.handlersPool, func(ctx context.Context) error {
        return fn(ctx, inputChannel, outputChannels)
    })
    conveyer.mutex.Unlock()
}

func (conveyer *Conveyer) Run(ctx context.Context) error {
    conveyer.mutex.Lock()
    if conveyer.running {
        conveyer.mutex.Unlock()
        return errors.New("conveyer already running")
    }
    conveyer.ctx, conveyer.cancel = context.WithCancel(ctx)
    conveyer.running = true
    conveyer.mutex.Unlock()

    defer func() {
        conveyer.mutex.Lock()
        conveyer.running = false
        conveyer.cancel()
        conveyer.closeAllChannels()
        conveyer.mutex.Unlock()
    }()

    if len(conveyer.handlersPool) == 0 {
        return nil
    }

    errorChannel := make(chan error, len(conveyer.handlersPool))
    var waitGroup sync.WaitGroup

    for _, handler := range conveyer.handlersPool {
        waitGroup.Add(1)
        go func(currentHandler func(context.Context) error) {
            defer waitGroup.Done()
            if handlerError := currentHandler(conveyer.ctx); handlerError != nil {
                select {
                case errorChannel <- handlerError:
                default:
                }
            }
        }(handler)
    }

    go func() {
        waitGroup.Wait()
        close(errorChannel)
    }()

    select {
    case <-conveyer.ctx.Done():
        waitGroup.Wait()
        return conveyer.ctx.Err()
    case handlerError, channelOpen := <-errorChannel:
        if channelOpen && handlerError != nil {
            conveyer.cancel()
            waitGroup.Wait()
            return handlerError
        }
        return nil
    }
}

func (conveyer *Conveyer) closeAllChannels() {
    for channelName, channel := range conveyer.channels {
        close(channel)
        delete(conveyer.channels, channelName)
    }
}

func (conveyer *Conveyer) Send(input string, data string) error {
    conveyer.mutex.RLock()
    channel, channelExists := conveyer.channels[input]
    conveyer.mutex.RUnlock()
    
    if !channelExists {
        return ErrChanNotFound
    }
    
    channel <- data
    return nil
}

func (conveyer *Conveyer) Recv(output string) (string, error) {
    conveyer.mutex.RLock()
    channel, channelExists := conveyer.channels[output]
    conveyer.mutex.RUnlock()
    
    if !channelExists {
        return "", ErrChanNotFound
    }
    
    data, channelOpen := <-channel
    if !channelOpen {
        return UndefinedData, nil
    }
    
    return data, nil
}
