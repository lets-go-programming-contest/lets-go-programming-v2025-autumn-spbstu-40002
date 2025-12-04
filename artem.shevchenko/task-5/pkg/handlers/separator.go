package handlers

import (
    "context"
)

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
    if len(outputs) == 0 {
        return nil
    }
    
    chanIndex := 0
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case data, channelOpen := <-input:
            if !channelOpen {
                return nil
            }
            
            outputChannel := outputs[chanIndex]
            select {
            case <-ctx.Done():
                return ctx.Err()
            case outputChannel <- data:
            }
            
            chanIndex = (chanIndex + 1) % len(outputs)
        }
    }
}