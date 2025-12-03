package handlers

import (
	"context"
	"fmt"
)

type Separator func(
	context.Context,
	chan string,
	[]chan string,
) error

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	cnt := 0

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context done: %w", ctx.Err())

		case val, ok := <-input:
			if !ok {
				return nil
			}

			out := outputs[cnt%len(outputs)]

			select {
			case <-ctx.Done():
				return fmt.Errorf("context done: %w", ctx.Err())

			case out <- val:
			}

			cnt++
		}
	}
}
