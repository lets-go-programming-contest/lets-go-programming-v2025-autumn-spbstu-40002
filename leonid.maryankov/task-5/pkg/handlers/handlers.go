package handlers

import (
	"context"
	"errors"
	"strings"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case v, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(v, "no decorator") {
				return errors.New("can't be decorated")
			}

			if !strings.HasPrefix(v, "decorated: ") {
				v = "decorated: " + v
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- v:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	i := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case v, ok := <-input:
			if !ok {
				return nil
			}

			if len(outputs) == 0 {
				continue
			}

			out := outputs[i%len(outputs)]
			i++

			select {
			case <-ctx.Done():
				return nil
			case out <- v:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	open := len(inputs)
	closed := make([]bool, len(inputs))

	for open > 0 {
		for idx := range inputs {
			if closed[idx] {
				continue
			}

			select {
			case <-ctx.Done():
				return nil

			case v, ok := <-inputs[idx]:
				if !ok {
					closed[idx] = true
					open--
					continue
				}

				if strings.Contains(v, "no multiplexer") {
					continue
				}

				select {
				case <-ctx.Done():
					return nil
				case output <- v:
				}
			}
		}
	}

	return nil
}
