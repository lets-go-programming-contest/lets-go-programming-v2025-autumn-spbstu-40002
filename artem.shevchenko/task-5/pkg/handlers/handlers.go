package handlers

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/sync/errgroup"
)

var ErrCantBeDecorated = errors.New("can't be decorated")


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
	index := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}
			target := outputs[index%len(outputs)]
			select {
			case target <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
			index++
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	group, childCtx := errgroup.WithContext(ctx)

	for _, ch := range inputs {
		ch := ch
		group.Go(func() error {
			for {
				select {
				case <-childCtx.Done():
					return nil
				case data, ok := <-ch:
					if !ok {
						return nil
					}
					if strings.Contains(data, "no multiplexer") {
						continue
					}
					select {
					case output <- data:
					case <-childCtx.Done():
						return nil
					}
				}
			}
		})
	}

	return group.Wait()
}