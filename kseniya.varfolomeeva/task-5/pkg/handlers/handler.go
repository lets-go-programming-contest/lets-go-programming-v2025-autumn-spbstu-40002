package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

const (
	PrefixToAdd    = "decorated: "
	TriggerNoDeco  = "no decorator"
	TriggerNoMux   = "no multiplexer"
)

var (
	ErrCannotDecorate = errors.New("can't be decorated")
	ErrNoOutputs      = errors.New("outputs slice must not be empty")
)

func PrefixDecorator(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, TriggerNoDeco) {
				return ErrCannotDecorate
			}

			if !strings.HasPrefix(data, PrefixToAdd) {
				data = PrefixToAdd + data
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func Multiplexer(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return ErrNoOutputs
	}

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, inputChan := range inputs {
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

					if strings.Contains(data, TriggerNoMux) {
						continue
					}

					select {
					case output <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(inputChan)
	}

	wg.Wait()

	return nil
}

func Separator(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrNoOutputs
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

			selectedOutput := outputs[currentIndex%len(outputs)]
			currentIndex++

			select {
			case selectedOutput <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

