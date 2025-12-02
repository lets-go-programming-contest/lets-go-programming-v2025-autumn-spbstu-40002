package handlers

import (
	"context"
	"strings"
)

const noMultiplexerData = "no multiplexer"

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	for {
		if err := ctx.Err(); err != nil {
			return errCancelledContext
		}

		for _, channel := range inputs {
			select {
			case value, ok := <-channel:
				if !ok {
					continue
				}

				if strings.Contains(value, noMultiplexerData) {
					continue
				}
				select {
				case <-ctx.Done():
					return errCancelledContext
				case output <- value:
				}
			case <-ctx.Done():
				return errCancelledContext
			default:
				continue
			}
		}
	}
}
