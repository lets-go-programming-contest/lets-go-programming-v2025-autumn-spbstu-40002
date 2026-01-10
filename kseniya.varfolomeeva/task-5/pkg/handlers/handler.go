package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

const (
	prefixToAdd   = "decorated: "
	triggerNoDeco = "no decorator"
	triggerNoMult = "no multiplexer"
)

var (
	ErrCantDecorate = errors.New("can't be decorated")
	ErrEmptyOutputs = errors.New("empty outputs")
)

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, triggerNoDeco) {
				return ErrCantDecorate
			}

			if !strings.HasPrefix(data, prefixToAdd) {
				data = prefixToAdd + data
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return ErrEmptyOutputs
	}

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(inputs))

	for _, chanal := range inputs {
		go func(inputChan chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-inputChan:
					if !ok {
						return
					}

					if strings.Contains(data, triggerNoMult) {
						continue
					}

					select {
					case output <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(chanal)
	}

	waitGroup.Wait()

	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrEmptyOutputs
	}

	currentIndex := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			targetIndex := currentIndex % len(outputs)
			currentIndex++

			select {
			case outputs[targetIndex] <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
