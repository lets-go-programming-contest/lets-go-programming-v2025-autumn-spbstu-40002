package handlers

import (
	"context"
	"strings"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, success := <-input:
			if !success {
				return nil
			}

			if strings.Contains(data, NoDecoratorFound) {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(data, Prefix) {
				data = Prefix + data
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}
