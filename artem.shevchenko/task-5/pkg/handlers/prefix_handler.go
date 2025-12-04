package handlers

import (
    "context"
    "strings"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case data, success := <-input:
            if !success {
                // Входной канал закрыт
                return nil
            }
            
            if strings.Contains(data, NoDecoratorFound) {
                return CantBeDecorated
            }
            
            if !strings.HasPrefix(data, Prefix) {
                data = Prefix + data
            }
            
            select {
            case <-ctx.Done():
                return ctx.Err()
            case output <- data:
            }
        }
    }
}