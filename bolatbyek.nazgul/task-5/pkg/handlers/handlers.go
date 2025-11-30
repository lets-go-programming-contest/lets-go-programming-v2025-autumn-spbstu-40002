package handlers

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

const prefix = "decorated: "

// PrefixDecoratorFunc is a data modifier that adds prefix "decorated: <original data>"
// to input data, but only if this prefix hasn't been added before.
// If input data contains substring "по decorator", the handler must terminate
// and return an error containing substring "can't be decorated".
func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("prefix decorator context error: %w", ctx.Err())
		case data, okChannel := <-input:
			if !okChannel {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
			}

			select {
			case <-ctx.Done():
				return fmt.Errorf("prefix decorator context error: %w", ctx.Err())
			case output <- data:
			}
		}
	}
}

// SeparatorFunc is a separator that distributes data from input channel to output channels
// based on sequential reception number. For example, for two output channels:
// first value goes to first channel, second to second, third to first, etc.
func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		for {
			select {
			case <-ctx.Done():
				return fmt.Errorf("separator context error: %w", ctx.Err())
			case _, okChannel := <-input:
				if !okChannel {
					return nil
				}
			}
		}
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("separator context error: %w", ctx.Err())
		case data, okChannel := <-input:
			if !okChannel {
				return nil
			}

			outputChan := outputs[index%len(outputs)]
			index++

			select {
			case <-ctx.Done():
				return fmt.Errorf("separator context error: %w", ctx.Err())
			case outputChan <- data:
			}
		}
	}
}

// MultiplexerFunc is a multiplexer that receives data from input channels and combines it.
// If data contains substring "по multiplexer", the data should be skipped (filtered).
// Continues working if at least one channel is open.
func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputs))

	for _, inputChan := range inputs {
		goroutine := func() {
			defer waitGroup.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case data, okChannel := <-inputChan:
					if !okChannel {
						return
					}

					if strings.Contains(data, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- data:
					}
				}
			}
		}
		go goroutine()
	}

	waitGroup.Wait()

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("multiplexer context error: %w", err)
	}

	return nil
}
