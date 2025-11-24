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
