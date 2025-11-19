package handlers

import (
	"context"
	"errors"
	"reflect"
	"strings"
)

const (
	decoratorPrefix   = "decorated: "
	noDecoratorMsg    = "по decorator"
	multiplexerFilter = "по multiplexer"
)

// PrefixDecoratorFunc is a data modifier that adds prefix "decorated: <original data>"
// to input data, but only if this prefix hasn't been added before.
// If input data contains substring "по decorator", the handler must terminate
// and return an error containing substring "can't be decorated"
func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			// Check if data contains "по decorator"
			if strings.Contains(data, noDecoratorMsg) {
				return errors.New("can't be decorated")
			}

			// Check if prefix already exists
			if strings.HasPrefix(data, decoratorPrefix) {
				// Prefix already exists, send as is
				select {
				case <-ctx.Done():
					return ctx.Err()
				case output <- data:
				}
			} else {
				// Add prefix
				decorated := decoratorPrefix + data
				select {
				case <-ctx.Done():
					return ctx.Err()
				case output <- decorated:
				}
			}
		}
	}
}

// MultiplexerFunc is a multiplexer that receives data from input channels and combines it.
// If data contains substring "по multiplexer", the data should be skipped (filtered)
// Continues working if at least one channel is open
func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	// Build dynamic select cases for all input channels
	selectCases := make([]reflect.SelectCase, 0, len(inputs)+1)
	selectCases = append(selectCases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ctx.Done()),
	})
	for _, input := range inputs {
		selectCases = append(selectCases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(input),
		})
	}

	for {
		chosen, value, ok := reflect.Select(selectCases)
		if chosen == 0 {
			// Context cancelled
			return ctx.Err()
		}

		if !ok {
			// Channel closed, rebuild select cases without this channel
			newCases := make([]reflect.SelectCase, 0, len(selectCases))
			newCases = append(newCases, selectCases[0]) // Keep context case
			for i := 1; i < len(selectCases); i++ {
				if i != chosen {
					newCases = append(newCases, selectCases[i])
				}
			}
			selectCases = newCases

			// If no input channels left, return
			if len(selectCases) == 1 {
				return nil
			}
			continue
		}

		data := value.String()
		// Filter data containing "по multiplexer"
		if strings.Contains(data, multiplexerFilter) {
			continue
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case output <- data:
		}
	}
}

// SeparatorFunc is a separator that distributes data from input channel to output channels
// based on sequential reception number. For example, for two output channels:
// first value goes to first channel, second to second, third to first, etc.
func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	index := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			// Distribute to next output channel in round-robin fashion
			select {
			case <-ctx.Done():
				return ctx.Err()
			case outputs[index] <- data:
				index = (index + 1) % len(outputs)
			}
		}
	}
}
