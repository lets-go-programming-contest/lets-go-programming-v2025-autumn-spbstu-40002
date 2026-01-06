package handlers

import (
	"context"
	"errors"
	"strings"
)

var errCanNotBeDecorated = errors.New("can't be decorated")

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

			if strings.Contains(data, "no decorator") {
				return errCanNotBeDecorated
			}

			if !strings.HasPrefix(data, "decorated: ") {
				data = "decorated: " + data
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

	idx := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}
			select {
			case outputs[idx%len(outputs)] <- data:
			case <-ctx.Done():

				return nil
			}

			idx++
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	doneChan := make(chan struct{}, len(inputs))

	for _, inputChannel := range inputs {
		go processInput(ctx, inputChannel, output, doneChan)
	}

	for i := 0; i < len(inputs); i++ { //nolint:warnamelen
		select {
		case <-doneChan:
		case <-ctx.Done():
			return nil
		}
	}

	return nil
}

func processInput(ctx context.Context, inputChannel <-chan string, output chan<- string, doneChan chan<- struct{}) {
	defer func() {
		doneChan <- struct{}{}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-inputChannel:
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
}
