package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	errCantBeDecorated = errors.New("can't be decorated")
	errEmptyOutputs    = errors.New("empty output")
)

const (
	prefix        = "decorated: "
	noDecorator   = "no decorator"
	noMultiplexer = "no multiplexer"
)

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, noDecorator) {
				return errCantBeDecorated
			}

			result := data
			if !strings.HasPrefix(data, prefix) {
				result = prefix + data
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- result:
			}
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	defer (func() {
		for _, channel := range outputs {
			close(channel)
		}
	})()

	if len(outputs) == 0 {
		return errEmptyOutputs
	}

	index := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case outputs[index] <- data:
			}

			index++
			if index >= len(outputs) {
				index = 0
			}
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	workersCount := len(inputs)
	for idx := 0; idx < workersCount; idx++ {
		inputChan := inputs[idx]

		go func(in <-chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, channelOpen := <-in:
					if !channelOpen {
						return
					}

					if strings.Contains(data, noMultiplexer) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- data:
					}
				}
			}
		}(inputChan)
	}

	select {
	case <-ctx.Done():
		return nil
	case <-done:
		return nil
	}
}
