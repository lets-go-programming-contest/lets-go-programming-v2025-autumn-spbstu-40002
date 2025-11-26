package handlers

import (
	"context"
	"errors"
	"strings"
)

var errPrefixDecoratorFuncCantBeDecorated = errors.New("handlers.PrefixDecoratorFunc: can't be decorated")

func PrefixDecoratorFunc(
	context context.Context,
	input chan string,
	output chan string,
) error {
	defer close(output)

	for {
		select {
		case <-context.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return errPrefixDecoratorFuncCantBeDecorated
			}

			if strings.HasPrefix(data, "decorated: ") {
				output <- data

				continue
			}

			select {
			case <-context.Done():
				return nil

			case output <- "decorated: " + data:
			}
		}
	}
}

func MultiplexerFunc(
	context context.Context,
	inputs []chan string,
	output chan string,
) error {
	defer close(output)

	activeInputs := make([]chan string, len(inputs))
	copy(activeInputs, inputs)

	for len(activeInputs) > 0 {
		select {
		case <-context.Done():
			return nil

		default:
			for i := 0; i < len(activeInputs); i++ { //nolint:varnamelen
				select {
				case <-context.Done():
					return nil

				case data, ok := <-activeInputs[i]:
					if !ok {
						activeInputs = append(activeInputs[:i], activeInputs[i+1:]...)
						i--

						continue
					}

					if strings.Contains(data, "no multiplexer") {
						continue
					}

					select {
					case <-context.Done():
						return nil

					case output <- data:
					}
				default:
				}
			}
		}
	}

	return nil
}

func SeparatorFunc(
	context context.Context,
	input chan string,
	outputs []chan string,
) error {
	defer (func() {
		for _, output := range outputs {
			close(output)
		}
	})()

	for i := 0; ; { //nolint:varnamelen
		select {
		case <-context.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case <-context.Done():
				return nil

			case outputs[i%len(outputs)] <- data:
				i++
			}
		}
	}
}
