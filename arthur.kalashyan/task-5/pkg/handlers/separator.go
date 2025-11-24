package handlers

import "context"

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		for {
			select {
			case <-ctx.Done():
				return nil
			case _, ok := <-input:
				if !ok {
					return nil
				}
			}
		}
	}
	idx := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case v, ok := <-input:
			if !ok {
				return nil
			}
			out := outputs[idx%len(outputs)]
			idx++
			select {
			case <-ctx.Done():
				return nil
			case out <- v:
			}
		}
	}
}
