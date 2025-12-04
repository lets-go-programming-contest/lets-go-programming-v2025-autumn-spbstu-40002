package handlers

import (
	"strings"
	"context"
)


func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output) // Close the output channel when finished.

	for {
		select {
		case <- ctx.Done():
			return ctx.Err()
		case data, success := <- input:
			if !success {
				return nil
			}

			// Check if a string contains "no decorator".
			if strings.Contains(data, NoDecoratorFound) {
				return CantBeDecorated
			}
			
			// Check if a string contains "decorated: ".
			if !strings.HasPrefix(data, Prefix) {
				data = Prefix + data
			}

			// Send the processed string to output.
			select {
			case <- ctx.Done():
				return ctx.Err()
			case output <- data:
			}
		}
	}
}
