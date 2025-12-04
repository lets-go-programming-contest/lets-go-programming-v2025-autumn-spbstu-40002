package handlers

import (
	"context"
	"strings"
	"sync"
)

const (
	decoratorPrefix  = "decorated: "
	noDecoratorMsg   = "no decorator"
	noMultiplexerMsg = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, noDecoratorMsg) {
				return &DecorationError{"can't be decorated"}
			}

			if !strings.HasPrefix(data, decoratorPrefix) {
				data = decoratorPrefix + data
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, out := range outputs {
			close(out)
		}
	}()

	if len(outputs) == 0 {
		return nil
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
			case outputs[index] <- data:
				index = (index + 1) % len(outputs)
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	return runMultiplexer(ctx, inputs, output)
}

func runMultiplexer(ctx context.Context, inputs []chan string, output chan string) error {
	var waitGroup sync.WaitGroup

	mergedChannel := make(chan string, len(inputs))

	startReaders(ctx, inputs, mergedChannel, &waitGroup)

	go closeMergedChannel(&waitGroup, mergedChannel)

	return writeToOutput(ctx, mergedChannel, output)
}

func startReaders(ctx context.Context, inputs []chan string, merged chan string, wg *sync.WaitGroup) {
	for _, inputChan := range inputs {
		wg.Add(1)

		go readFromChannel(ctx, inputChan, merged, wg)
	}
}

func readFromChannel(ctx context.Context, inputChan <-chan string, merged chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-inputChan:
			if !ok {
				return
			}

			if strings.Contains(data, noMultiplexerMsg) {
				continue
			}

			select {
			case merged <- data:
			case <-ctx.Done():
				return
			}
		}
	}
}

func closeMergedChannel(wg *sync.WaitGroup, merged chan string) {
	wg.Wait()
	close(merged)
}

func writeToOutput(ctx context.Context, merged <-chan string, output chan<- string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-merged:
			if !ok {
				return nil
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

type DecorationError struct {
	Msg string
}

func (e *DecorationError) Error() string {
	return e.Msg
}
