package conveyer

import (
    "context"
    "errors"
    "sync"
)

type Conveyer struct {
    mutex           sync.RWMutex
    channelSize  int
    channels     map[string]chan string
    handlersPool []func(context.Context) error
    ctx          context.Context
    cancel       context.CancelFunc
    running      bool
}

func New(size int) *Conveyer {
    return &Conveyer{
        channelSize:  size,
        channels:     make(map[string]chan string),
        handlersPool: make([]func(context.Context) error, 0),
        running:      false,
    }
}

func (conv *Conveyer) getOrCreateChannel(name string) chan string {
    conv.mutex.Lock()
    defer conv.mutex.Unlock()
    
    if channel, exists := conv.channels[name]; exists {
        return channel
    }
    
    channel := make(chan string, conv.channelSize)
    conv.channels[name] = channel
    return channel
}

func (conv *Conveyer) RegisterDecorator(
    fn func(ctx context.Context, input chan string, output chan string) error,
    input string,
    output string,
) {
    inputChan := conv.getOrCreateChannel(input)
    outputChan := conv.getOrCreateChannel(output)
    
    conv.mutex.Lock()
    conv.handlersPool = append(conv.handlersPool, func(ctx context.Context) error {
        return fn(ctx, inputChan, outputChan)
    })
    conv.mutex.Unlock()
}

func (conv *Conveyer) RegisterMultiplexer(
    fn func(ctx context.Context, inputs []chan string, output chan string) error,
    inputs []string,
    output string,
) {
    inputChannels := make([]chan string, len(inputs))
    for i, input := range inputs {
        inputChannels[i] = conv.getOrCreateChannel(input)
    }
    outputChan := conv.getOrCreateChannel(output)
    
    conv.mutex.Lock()
    conv.handlersPool = append(conv.handlersPool, func(ctx context.Context) error {
        return fn(ctx, inputChannels, outputChan)
    })
    conv.mutex.Unlock()
}

func (conv *Conveyer) RegisterSeparator(
    fn func(ctx context.Context, input chan string, outputs []chan string) error,
    input string,
    outputs []string,
) {
    inputChan := conv.getOrCreateChannel(input)
    outputChannels := make([]chan string, len(outputs))
    for i, output := range outputs {
        outputChannels[i] = conv.getOrCreateChannel(output)
    }
    
    conv.mutex.Lock()
    conv.handlersPool = append(conv.handlersPool, func(ctx context.Context) error {
        return fn(ctx, inputChan, outputChannels)
    })
    conv.mutex.Unlock()
}

func (conv *Conveyer) Run(ctx context.Context) error {
    conv.mutex.Lock()
    if conv.running {
        conv.mutex.Unlock()
        return errors.New("conveyer already running")
    }
    conv.ctx, conv.cancel = context.WithCancel(ctx)
    conv.running = true
    conv.mutex.Unlock()
    
    defer func() {
        conv.mutex.Lock()
        conv.running = false
        conv.cancel()
        conv.mutex.Unlock()
    }()
    
    if len(conv.handlersPool) == 0 {
        return nil
    }
    
    errCh := make(chan error, len(conv.handlersPool))
    var waitg sync.WaitGroup
    
    for _, handler := range conv.handlersPool {
        waitg.Add(1)
        go func(h func(context.Context) error) {
            defer waitg.Done()
            if err := h(conv.ctx); err != nil {
                select {
                case errCh <- err:
                default:
                }
            }
        }(handler)
    }
    
    go func() {
        waitg.Wait()
        close(errCh)
    }()
    
    select {
    case <-conv.ctx.Done():
        return conv.ctx.Err()
    case err := <-errCh:
        conv.cancel()
        return err
    }
}

func (conv *Conveyer) Send(input string, data string) error {
    conv.mutex.RLock()
    channel, exists := conv.channels[input]
    conv.mutex.RUnlock()
    
    if !exists {
        return ErrChanNotFound
    }
    
    channel <- data
    return nil
}

func (conv *Conveyer) Recv(output string) (string, error) {
    conv.mutex.RLock()
    channel, exists := conv.channels[output]
    conv.mutex.RUnlock()
    
    if !exists {
        return "", ErrChanNotFound
    }
    
    data, succes := <-channel
    if !succes {
        return UndefinedData, nil
    }
    
    return data, nil
}
