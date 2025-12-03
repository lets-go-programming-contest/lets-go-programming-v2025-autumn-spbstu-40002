package handlers

import (
	"context"
	"errors"
	"strings"
)

var errPrefixDecorator = errors.New("handlers.PrefixDecoratorFunc: can't be decorated")

type Decorator func(
	context.Context,
	chan string,
	chan string,
) error

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case val, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(val, "no decorator") {
				return errPrefixDecorator
			}
			if !strings.HasPrefix(val, "decorated: ") {
				val = "decorated:" + val
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- val:
			}

		}
	}
}
