package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case value, open := <-input:
			if !open {
				return nil
			}

			if strings.Contains(value, "no decorator") {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(value, "decorated: ") {
				value = "decorated: " + value
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- value:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case value, open := <-input:
			if !open {
				return nil
			}

			if len(outputs) == 0 {
				continue
			}

			outCh := outputs[index%len(outputs)]
			index++

			select {
			case <-ctx.Done():
				return nil
			case outCh <- value:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputs))

	for _, ch := range inputs {
		inputCh := ch

		go func(inputCh chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case value, open := <-inputCh:
					if !open {
						return
					}

					if strings.Contains(value, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- value:
					}
				}
			}
		}(inputCh)
	}

	waitGroup.Wait()

	return nil
}
