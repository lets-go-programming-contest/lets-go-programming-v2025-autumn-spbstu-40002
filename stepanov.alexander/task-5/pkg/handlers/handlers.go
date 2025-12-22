package handlers

import (
	"context"
)

func PrefixDecorator(ctx context.Context, in, out chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-in:
			if !ok {
				return nil
			}
			select {
			case out <- "prefix: " + data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func Multiplexer(ctx context.Context, inputs []chan string, output chan string) error {
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
					select {
					case output <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(input)
	}
	return nil
}

func Separator(ctx context.Context, input chan string, outputs []chan string) error {
	idx := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}
			outputs[idx] <- data
			idx = (idx + 1) % len(outputs)
		}
	}
}
