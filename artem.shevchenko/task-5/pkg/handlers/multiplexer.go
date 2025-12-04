package handlers

import (
    "context"
    "strings"
    "sync"
)

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
    var waitGroup sync.WaitGroup
    doneChannel := make(chan struct{})
    
    reader := func(inputChannel chan string) {
        defer waitGroup.Done()
        for {
            select {
            case <-ctx.Done():
                return
            case <-doneChannel:
                return
            case data, channelOpen := <-inputChannel:
                if !channelOpen {
                    return
                }
                if strings.Contains(data, NoMultiplexer) {
                    continue
                }
                select {
                case <-ctx.Done():
                    return
                case <-doneChannel:
                    return
                case output <- data:
                }
            }
        }
    }
    
    for _, inputChannel := range inputs {
        waitGroup.Add(1)
        go reader(inputChannel)
    }
    
    waitGroup.Wait()
    close(doneChannel)
    return nil
}
