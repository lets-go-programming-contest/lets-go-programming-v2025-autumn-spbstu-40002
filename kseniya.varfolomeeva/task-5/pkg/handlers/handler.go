package handlers

import (
    "context"
    "errors"
    "strings"
    "sync"
)

const (
    TagToPrepend = "decorated: "
    NoTagSignal  = "no tag"
    SkipMixSignal = "skip mix"
)

var (
    ErrTagRejected = errors.New("tag rejected")
    ErrEmptyTargets = errors.New("empty targets")
)

func ApplyTagFunc(ctx context.Context, input, output chan string) error {
    for {
        select {
        case <-ctx.Done():
            return nil
        case item, active := <-input:
            if !active {
                return nil
            }

            if strings.Contains(item, NoTagSignal) {
                return ErrTagRejected
            }

            if !strings.HasPrefix(item, TagToPrepend) {
                item = TagToPrepend + item
            }

            select {
            case output <- item:
            case <-ctx.Done():
                return nil
            }
        }
    }
}

func MixStreamsFunc(ctx context.Context, inputs []chan string, output chan string) error {
    if len(inputs) == 0 {
        return ErrEmptyTargets
    }

    var mixer sync.WaitGroup
    mixer.Add(len(inputs))

    for _, source := range inputs {
        go func(inChan chan string) {
            defer mixer.Done()

            for {
                select {
                case <-ctx.Done():
                    return
                case content, active := <-inChan:
                    if !active {
                        return
                    }

                    if strings.Contains(content, SkipMixSignal) {
                        continue
                    }

                    select {
                    case output <- content:
                    case <-ctx.Done():
                        return
                    }
                }
            }
        }(source)
    }

    mixer.Wait()
    return nil
}

func RouteStreamFunc(ctx context.Context, input chan string, outputs []chan string) error {
    if len(outputs) == 0 {
        return ErrEmptyTargets
    }

    counter := 0

    for {
        select {
        case <-ctx.Done():
            return nil
        case element, active := <-input:
            if !active {
                return nil
            }

            route := outputs[counter%len(outputs)]
            counter++

            select {
            case route <- element:
            case <-ctx.Done():
                return nil
            }
        }
    }
}
