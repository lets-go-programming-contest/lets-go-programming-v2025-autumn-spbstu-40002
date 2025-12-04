package handlers

import (
    "context"
)

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
    defer func() {
        // Close all output channels when finished.
        for _, out := range outputs {
            close(out)
        }
    }()
    
    if len(outputs) == 0 {
        return nil
    }
    
    chanIndex := 0
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case data, success := <- input:
            if !success {
                return nil
            }
            
            // Select the output channel in order.
            out := outputs[chanIndex]
            // Send data to the selected channel.
            select {
            case <-ctx.Done():
                return ctx.Err()
            case out <- data:
            }
            // Move on to the next channel.
            chanIndex = (chanIndex + 1) % len(outputs)
        }
    }
}
