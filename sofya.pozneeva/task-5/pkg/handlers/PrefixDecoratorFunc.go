package handlers

import (
	"context"
	"errors"
	"strings"
)

var errStringNotDecorated = errors.New("can’t be decorated")
var errStringDecorated = errors.New("string was decorated earlier")

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return errCancelledContext

		case inputString, ok := <-input:
			if !ok {
				return errInvalidChan
			}

			if strings.Contains(inputString, "decorator:") {
				return errStringDecorated
			}

			if strings.Contains(inputString, "no decorator") {
				return errStringNotDecorated
			}

			decorated := "decorator:" + inputString

			select {
			case output <- decorated:
				return nil // Успешно обработали одно значение
			case <-ctx.Done():
				return errCancelledContext
			}
		}
	}
}
