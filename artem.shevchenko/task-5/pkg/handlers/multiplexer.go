package handlers

import (
    "context"
    "strings"
    "sync"
)

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
    var waitg sync.WaitGroup

    // For each input channel start a goroutine.
    for _, input := range inputs {
        waitg.Add(1)
        go func(inChan chan string) {
            defer waitg.Done()
            for {
                select {
                case <- ctx.Done():
                    return // Context canceled.
                case data, success := <- inChan:
                    if !success {
                        return // Channel closed.
                    }
                    // Skipping data with "no multiplexer".
                    if strings.Contains(data, NoMultiplexer) {
                        continue
                    }
                    // send to output
                    select {
                    case <- ctx.Done():
                        return
                    case output <- data:
                    }
                }
            }
        }(input)
    }

    // Waiting for all goroutines to complete.
    go func() {
        waitg.Wait()
        close(output) // Close the output channel after all goroutines have completed.
    }()

    // Expect the context to be cancelled.
    <- ctx.Done()
    return ctx.Err()
}
