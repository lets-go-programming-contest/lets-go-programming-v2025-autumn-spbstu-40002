package handlers

import (
	"context"
	"fmt"
	"strings"
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
				return fmt.Errorf("can't be decorated")
			}

			if !strings.HasPrefix(data, "decorated: ") {
				data = "decorated: " + data
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- data:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	counter := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			index := counter % len(outputs)
			counter++

			select {
			case <-ctx.Done():
				return ctx.Err()
			case outputs[index] <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	combined := make(chan string)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, input := range inputs {
		go func(in chan string) {
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-in:
					if !ok {
						return
					}

					if strings.Contains(data, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case combined <- data:
					}
				}
			}
		}(input)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-combined:
			if !ok {
				close(output)
				return nil
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- data:
			}
		}
	}
}
