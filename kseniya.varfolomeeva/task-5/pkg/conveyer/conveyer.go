package conveyer

import (
    "context"
    "errors"
    "fmt"
    "sync"

    "golang.org/x/sync/errgroup"
)

var ErrChanMissing = errors.New("channel missing")

const UndefinedData = "undefined"

type ConveyerInterface interface {
    AddHandler(handler func(context.Context, chan string, chan string) error, inName, outName string)
    AddCombiner(handler func(context.Context, []chan string, chan string) error, inNames []string, outName string)
    AddDivider(handler func(context.Context, chan string, []chan string) error, inName string, outNames []string)
    Execute(ctx context.Context) error
    Insert(inName string, data string) error
    Extract(outName string) (string, error)
}

type operation func(context.Context) error

type StreamProcessor struct {
    capacity   int
    streams    map[string]chan string
    operations []operation
    lock       sync.RWMutex
}

func Build(size int) *StreamProcessor {
    return &StreamProcessor{
        capacity:   size,
        streams:    make(map[string]chan string),
        operations: make([]operation, 0),
        lock:       sync.RWMutex{},
    }
}

func (sp *StreamProcessor) getStream(name string) chan string {
    sp.lock.Lock()
    defer sp.lock.Unlock()

    if stream, present := sp.streams[name]; present {
        return stream
    }

    newStream := make(chan string, sp.capacity)
    sp.streams[name] = newStream
    return newStream
}

func (sp *StreamProcessor) findStream(name string) (chan string, error) {
    sp.lock.RLock()
    defer sp.lock.RUnlock()

    if stream, present := sp.streams[name]; present {
        return stream, nil
    }

    return nil, ErrChanMissing
}

func (sp *StreamProcessor) AddHandler(
    handler func(context.Context, chan string, chan string) error,
    input, output string,
) {
    inStream := sp.getStream(input)
    outStream := sp.getStream(output)

    op := func(ctx context.Context) error {
        return handler(ctx, inStream, outStream)
    }

    sp.addOperation(op)
}

func (sp *StreamProcessor) AddCombiner(
    handler func(context.Context, []chan string, chan string) error,
    inputs []string,
    output string,
) {
    inStreams := make([]chan string, len(inputs))
    for i, name := range inputs {
        inStreams[i] = sp.getStream(name)
    }

    outStream := sp.getStream(output)

    op := func(ctx context.Context) error {
        return handler(ctx, inStreams, outStream)
    }

    sp.addOperation(op)
}

func (sp *StreamProcessor) AddDivider(
    handler func(context.Context, chan string, []chan string) error,
    input string,
    outputs []string,
) {
    inStream := sp.getStream(input)

    outStreams := make([]chan string, len(outputs))
    for i, name := range outputs {
        outStreams[i] = sp.getStream(name)
    }

    op := func(ctx context.Context) error {
        return handler(ctx, inStream, outStreams)
    }

    sp.addOperation(op)
}

func (sp *StreamProcessor) Insert(inName string, data string) error {
    stream, err := sp.findStream(inName)
    if err != nil {
        return err
    }

    select {
    case stream <- data:
        return nil
    default:
        return fmt.Errorf("stream %s overflow", inName)
    }
}

func (sp *StreamProcessor) Extract(outName string) (string, error) {
    stream, err := sp.findStream(outName)
    if err != nil {
        return "", err
    }

    value, active := <-stream
    if !active {
        return UndefinedData, nil
    }

    return value, nil
}

func (sp *StreamProcessor) shutdownStreams() {
    sp.lock.Lock()
    defer sp.lock.Unlock()

    for _, stream := range sp.streams {
        close(stream)
    }
}

func (sp *StreamProcessor) Execute(ctx context.Context) error {
    sp.lock.RLock()
    defer sp.lock.RUnlock()

    group, ctx := errgroup.WithContext(ctx)

    for _, op := range sp.operations {
        currentOp := op
        group.Go(func() error {
            return currentOp(ctx)
        })
    }

    err := group.Wait()
    sp.shutdownStreams()

    if err != nil {
        return fmt.Errorf("stream processing failed: %w", err)
    }

    return nil
}

func (sp *StreamProcessor) addOperation(op operation) {
    sp.lock.Lock()
    defer sp.lock.Unlock()
    sp.operations = append(sp.operations, op)
}
