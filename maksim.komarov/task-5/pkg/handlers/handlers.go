package handlers

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		case v, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(v, "no decorator") {
				return fmt.Errorf("can't be decorated")
			}

			output <- "decorated: " + v
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
		case v, ok := <-input:
			if !ok {
				return nil
			}

			outputs[index] <- v
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
				case v, ok := <-inCh:
					if !ok {
						return
					}
					if strings.Contains(v, "no multiplexer") {
						continue
					}
					output <- v
				}
			}
		}()
	}

	waitGroup.Wait()
	return nil
}
