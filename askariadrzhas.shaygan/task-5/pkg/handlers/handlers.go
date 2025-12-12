package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	errCantBeDecorated = errors.New("cannot apply decoration")
	errInputsEmpty     = errors.New("no input channels")
	errOutputsEmpty    = errors.New("no output channels")
)

const (
	noDecoratorSub   = "no decorator"
	noMultiplexerSub = "no multiplexer"
	decoratedSub     = "decorated: "
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case data, active := <-input:
			if !active {
				return nil
			}

			if strings.Contains(data, noDecoratorSub) {
				return errCantBeDecorated
			}

			result := data
			if !strings.HasPrefix(data, decoratedSub) {
				result = decoratedSub + data
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- result:
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	count := len(outputs)

	if count == 0 {
		return errOutputsEmpty
	}

	position := 0

	for {
		select {
		case data, active := <-input:
			if !active {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case outputs[position] <- data:
			}

			position = (position + 1) % count
		case <-ctx.Done():
			return nil
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup

	inputCount := len(inputs)
	if inputCount == 0 {
		return errInputsEmpty
	}

	wg.Add(inputCount)

	for _, inp := range inputs {
		go func(ch chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, active := <-ch:
					if !active {
						return
					}

					if !strings.Contains(data, noMultiplexerSub) {
						select {
						case <-ctx.Done():
							return
						case output <- data:
						}
					}
				}
			}
		}(inp)
	}

	wg.Wait()

	return nil
}
