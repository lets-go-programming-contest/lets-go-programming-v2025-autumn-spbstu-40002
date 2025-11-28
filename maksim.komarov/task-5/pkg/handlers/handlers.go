package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	const prefix = "decorated: "

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

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
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

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputs))

	for _, channelInstance := range inputs {
		inputChan := channelInstance

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
