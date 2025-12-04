package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrCantBeDecorated = errors.New("can’t be decorated")

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	const prefix = "decorated: "
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(data, "no decorator") {
				return ErrCantBeDecorated
			}
			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
			}
			select {
			case output <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}


func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}
	idx := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}
			target := outputs[idx%len(outputs)]
			select {
			case target <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
			idx++
		}
	}
}


func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()


	for _, ch := range inputs {
		ch := ch
		go func() {
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

	<-ctx.Done()
	return ctx.Err()
}