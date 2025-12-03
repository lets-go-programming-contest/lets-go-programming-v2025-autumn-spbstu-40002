package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	errInputs      = errors.New("inputs can't be empty")
	errOutputs     = errors.New("outputs can't be empty")
	errNoDecorator = errors.New("can't be decorated")
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return errNoDecorator
			}

			res := data
			if !strings.HasPrefix(data, "decorated: ") {
				res = "decorated: " + data
			}
			select {
			case output <- res:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return errOutputs
	}

	index := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[index] <- data:
			case <-ctx.Done():
				return ctx.Err()
			}

			index = (index + 1) % len(outputs)
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return errInputs
	}

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, input := range inputs {
		go func(ch <-chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(data, "no multiplexer") {
						continue
					}

					select {
					case output <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(input)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(output)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
