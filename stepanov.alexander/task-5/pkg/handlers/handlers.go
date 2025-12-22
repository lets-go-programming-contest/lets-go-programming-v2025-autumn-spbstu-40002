package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrNoDecorator = errors.New("can't be decorated")

const decoratedPrefix = "decorated: "

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrNoDecorator
			}

			if !strings.HasPrefix(data, decoratedPrefix) {
				data = decoratedPrefix + data
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup

	wg.Add(len(inputs))
	for _, inputChan := range inputs {
		ch := inputChan

		go func() {
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
		}()
	}

	wg.Wait()

	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return errors.New("empty outputs")
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[index] <- data:
			case <-ctx.Done():
				return nil
			}

			index = (index + 1) % len(outputs)
		}
	}
}
