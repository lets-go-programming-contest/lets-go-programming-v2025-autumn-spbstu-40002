package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	const prefix = "decorated: "
	const errorSubstr = "can't be decorated"

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return errors.New(errorSubstr)
			}

			prefixedData := data
			if !strings.HasPrefix(data, prefix) {
				prefixedData = prefix + data
			}

			select {
			case output <- prefixedData:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	var wg sync.WaitGroup

	for i := range inputs {
		wg.Add(1)
		go func(input chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-input:
					if !ok {
						return
					}

					if strings.Contains(data, "no multiplexer") {
						continue
					}

					select {
					case output <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(inputs[i])
	}

	wg.Wait()
	return nil
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	if len(outputs) == 0 {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case _, ok := <-input:
				if !ok {
					return nil
				}
			}
		}
	}

	outputIndex := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			currentOutput := outputs[outputIndex]

			select {
			case currentOutput <- data:
				outputIndex = (outputIndex + 1) % len(outputs)
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
