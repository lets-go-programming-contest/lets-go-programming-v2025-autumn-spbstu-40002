package handlers

import (
    "context"
    "strings"
    "sync"
)

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
    var wg sync.WaitGroup
    done := make(chan struct{})
    
    reader := func(inChan chan string) {
        defer wg.Done()
        for {
            select {
            case <-ctx.Done():
                return
            case <-done:
                return
            case data, ok := <- inChan:
                if !ok {
                    return
                }
                if strings.Contains(data, NoMultiplexer) {
                    continue
                }
                select {
                case <-ctx.Done():
                    return
                case <-done:
                    return
                case output <- data:
                }
            }
        }
    }
    
    for _, inChan := range inputs {
        wg.Add(1)
        go reader(inChan)
    }
    
    <-ctx.Done()
    close(done)
    
    wg.Wait()
    return ctx.Err()
}
