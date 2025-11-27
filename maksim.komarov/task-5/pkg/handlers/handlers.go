package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		case value, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(value, "no decorator") {
				return ErrCantBeDecorated
			}
			output <- "decorated: " + value
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}
	index := 0
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		case value, ok := <-input:
			if !ok {
				return nil
			}
			outputs[index] <- value
			index++
			if index == len(outputs) {
				index = 0
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(inputs))
	for _, ch := range inputs {
		inCh := ch
		go func() {
			defer waitGroup.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case value, ok := <-inCh:
					if !ok {
						return
					}
					if strings.Contains(value, "no multiplexer") {
						continue
					}
					output <- value
				}
			}
		}()
	}
	waitGroup.Wait()
	return nil
}
