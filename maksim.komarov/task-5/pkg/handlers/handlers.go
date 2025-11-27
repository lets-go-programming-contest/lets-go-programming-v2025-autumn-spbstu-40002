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

		case v, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(v, "no decorator") {
				return ErrCantBeDecorated
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
	var wg sync.WaitGroup

	wg.Add(len(inputs))

	for _, ch := range inputs {
		in := ch

		runner := func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case v, ok := <-in:
					if !ok {
						return
					}

					if strings.Contains(v, "no multiplexer") {
						continue
					}

					output <- v
				}
			}
		}

		go runner()
	}

	wg.Wait()

	return nil
}
