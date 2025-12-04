package conveyer

import (
    "context"
    
    "golang.org/x/sync/errgroup"
)

type Conveyer struct {
    channelSize  int
    channels     map[string]chan string
    handlersPool []func(context.Context) error
}

func New(channelSize int) *Conveyer {
    return &Conveyer{
        channelSize:  channelSize,
        channels:     make(map[string]chan string),
        handlersPool: make([]func(context.Context) error, 0),
    }
}

func (conveyer *Conveyer) makeChannels(names ...string) {
    for _, name := range names {
        if _, channelExists := conveyer.channels[name]; !channelExists {
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
    conveyer.handlersPool = append(conveyer.handlersPool, func(ctx context.Context) error {
        return fn(ctx, conveyer.channels[input], conveyer.channels[output])
    })
}

func (conveyer *Conveyer) RegisterMultiplexer(
    fn func(ctx context.Context, inputs []chan string, output chan string) error,
    inputs []string,
    output string,
) {
    conveyer.makeChannels(inputs...)
    conveyer.makeChannels(output)
    conveyer.handlersPool = append(conveyer.handlersPool, func(ctx context.Context) error {
        inputChannels := make([]chan string, len(inputs))
        for index, input := range inputs {
            inputChannels[index] = conveyer.channels[input]
        }
        return fn(ctx, inputChannels, conveyer.channels[output])
    })
}

func (conveyer *Conveyer) RegisterSeparator(
    fn func(ctx context.Context, input chan string, outputs []chan string) error,
    input string,
    outputs []string,
) {
    conveyer.makeChannels(input)
    conveyer.makeChannels(outputs...)
    conveyer.handlersPool = append(conveyer.handlersPool, func(ctx context.Context) error {
        outputChannels := make([]chan string, len(outputs))
        for index, output := range outputs {
            outputChannels[index] = conveyer.channels[output]
        }
        return fn(ctx, conveyer.channels[input], outputChannels)
    })
}

func (conveyer *Conveyer) Run(ctx context.Context) error {
    handlersGroup, handlersContext := errgroup.WithContext(ctx)
    
    for _, handler := range conveyer.handlersPool {
        currentHandler := handler
        handlersGroup.Go(func() error {
            return currentHandler(handlersContext)
        })
    }
    
    return handlersGroup.Wait()
}

func (conveyer *Conveyer) Send(input string, data string) error {
    channel, channelExists := conveyer.channels[input]
    if !channelExists {
        return ErrChanNotFound
    }
    
    channel <- data
    return nil
}

func (conveyer *Conveyer) Recv(output string) (string, error) {
    channel, channelExists := conveyer.channels[output]
    if !channelExists {
        return "", ErrChanNotFound
    }
    
    data, channelOpen := <-channel
    if !channelOpen {
        return UndefinedData, nil
    }
    
    return data, nil
}
