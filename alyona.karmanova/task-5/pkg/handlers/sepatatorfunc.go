package handlers

import "context"

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
			return ctx.Err()

		case val, ok := <-input:
			if !ok {
				return nil
			}

			out := outputs[cnt%len(outputs)]

			select {
			case <-ctx.Done():
				return ctx.Err()

			case out <- val:
			}

			cnt++
		}
	}
}
