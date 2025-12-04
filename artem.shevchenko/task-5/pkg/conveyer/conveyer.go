package conveyer

import (
    "context"
    "errors"
    "sync"

    "golang.org/x/sync/errgroup"
)

type Conveyer struct {
    mutex        sync.RWMutex
    channelSize  int
    channels     map[string]chan string
    handlersPool []func(context.Context) error
    ctx          context.Context
    cancel       context.CancelFunc
}

func New(size int) *Conveyer {
    return &Conveyer{
        channelSize:  size,
        channels:     make(map[string]chan string),
        handlersPool: make([]func(context.Context) error, 0),
    }
}

func (conveyer *Conveyer) makeChannels(names ...string) {
    conveyer.mutex.Lock()
    defer conveyer.mutex.Unlock()
    
    for _, name := range names {
        if _, exists := conveyer.channels[name]; !exists {
            conveyer.channels[name] = make(chan string, conveyer.channelSize)
        }
    }
}

func (conveyer *Conveyer) RegisterDecorator(
    fn func(ctx context.Context, input chan string, output chan string) error,
    input string,
    output string,
) {
    conveyer.makeChannels(input, output)
    
    conveyer.mutex.Lock()
    conveyer.handlersPool = append(conveyer.handlersPool, func(ctx context.Context) error {
        return fn(ctx, conveyer.channels[input], conveyer.channels[output])
    })
    conveyer.mutex.Unlock()
}

func (conveyer *Conveyer) RegisterMultiplexer(
    fn func(ctx context.Context, inputs []chan string, output chan string) error,
    inputs []string,
    output string,
) {
    conveyer.makeChannels(inputs...)
    conveyer.makeChannels(output)
    
    conveyer.mutex.Lock()
    conveyer.handlersPool = append(conveyer.handlersPool, func(ctx context.Context) error {
        inputChannels := make([]chan string, len(inputs))
        for index, input := range inputs {
            inputChannels[index] = conveyer.channels[input]
        }
        return fn(ctx, inputChannels, conveyer.channels[output])
    })
    conveyer.mutex.Unlock()
}

func (conveyer *Conveyer) RegisterSeparator(
    fn func(ctx context.Context, input chan string, outputs []chan string) error,
    input string,
    outputs []string,
) {
    conveyer.makeChannels(input)
    conveyer.makeChannels(outputs...)
    
    conveyer.mutex.Lock()
    conveyer.handlersPool = append(conveyer.handlersPool, func(ctx context.Context) error {
        outputChannels := make([]chan string, len(outputs))
        for index, output := range outputs {
            outputChannels[index] = conveyer.channels[output]
        }
        return fn(ctx, conveyer.channels[input], outputChannels)
    })
    conveyer.mutex.Unlock()
}

func (conveyer *Conveyer) Run(ctx context.Context) error {
    conveyer.mutex.Lock()
    conveyer.ctx, conveyer.cancel = context.WithCancel(ctx)
    defer conveyer.cancel()
    conveyer.mutex.Unlock()
    
    defer conveyer.closeAllChannels()
    
    if len(conveyer.handlersPool) == 0 {
        return nil
    }
    
    handlersGroup, handlersContext := errgroup.WithContext(conveyer.ctx)
    
    for _, handler := range conveyer.handlersPool {
        currentHandler := handler
        handlersGroup.Go(func() error {
            return currentHandler(handlersContext)
        })
    }
    
    return handlersGroup.Wait()
}

func (conveyer *Conveyer) closeAllChannels() {
    conveyer.mutex.Lock()
    defer conveyer.mutex.Unlock()
    
    for name, channel := range conveyer.channels {
        select {
        case <-channel:
            // Канал уже закрыт или пуст
        default:
            close(channel)
        }
        delete(conveyer.channels, name)
    }
}

func (conveyer *Conveyer) Send(input string, data string) error {
    conveyer.mutex.RLock()
    channel, channelExists := conveyer.channels[input]
    ctx := conveyer.ctx
    conveyer.mutex.RUnlock()
    
    if !channelExists {
        return ErrChanNotFound
    }
    
    if ctx == nil {
        return ErrConvNotRun
    }
    
    select {
    case <-ctx.Done():
        return ctx.Err()
    case channel <- data:
        return nil
    }
}

func (conveyer *Conveyer) Recv(output string) (string, error) {
    conveyer.mutex.RLock()
    channel, channelExists := conveyer.channels[output]
    ctx := conveyer.ctx
    conveyer.mutex.RUnlock()
    
    if !channelExists {
        return "", ErrChanNotFound
    }
    
    if ctx == nil {
        return "", ErrConvNotRun
    }
    
    select {
    case <-ctx.Done():
        return "", ctx.Err()
    case data, channelOpen := <-channel:
        if !channelOpen {
            return UndefinedData, nil
        }
        return data, nil
    }
}