package handlers

import (
	"context"
	"errors"
)

var errInvalidChan = errors.New("don't read value from chan")
var errCancelledContext = errors.New("context is cancelled")

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	if ctx.Err() != nil {
		return errCancelledContext
	}
	for _, channel := range outputs {
		select {
		case value, ok := <-input:
			if !ok {
				return errInvalidChan
			}
			select {
			case channel <- value:
				continue
			case <-ctx.Done():
				return errCancelledContext
			}

		case <-ctx.Done():
			return errCancelledContext
		}
	}

	return nil
}
