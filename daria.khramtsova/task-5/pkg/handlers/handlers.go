package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	errCantDecorate = errors.New("can't be decorated")
	errEmptyInputs  = errors.New("input is empty")
	errEmptyOutputs = errors.New("output is empty")
)

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case value, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(value, "no decorator") {
				return errCantDecorate
			}

			if !strings.HasPrefix(value, "decorated: ") {
				value = "decorated: " + value
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- value:
			}
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	if len(outputs) == 0 {
		return errEmptyOutputs
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case value, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case outputs[index] <- value:
			}

			index = (index + 1) % len(outputs)
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	if len(inputs) == 0 {
		return errEmptyInputs
	}

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, ch := range inputs {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case value, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(value, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- value:
					}
				}
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil
	case <-done:
		return nil
	}
}
