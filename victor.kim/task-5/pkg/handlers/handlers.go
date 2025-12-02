package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

const (
	Prefix               = "decorated: "
	NoDecoratorMessage   = "no decorator"
	NoMultiplexerMessage = "no multiplexer"
	Undefined            = "undefined"
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

			if strings.Contains(data, NoDecoratorMessage) {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(data, Prefix) {
				data = Prefix + data
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
		for _, output := range outputs {
			close(output)
		}
	}()

	if len(outputs) == 0 {
		return nil
	}

	counter := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			outputIndex := counter % len(outputs)
			counter++

			select {
			case outputs[outputIndex] <- data:
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

	waitGroup := sync.WaitGroup{}

	for _, inputChannel := range inputs {
		waitGroup.Add(1)

		channel := inputChannel
		readFunc := func() {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-channel:
					if !ok {
						return
					}

					if strings.Contains(data, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- data:
					}
				}
			}
		}

		go readFunc()
	}

	waitGroup.Wait()

	return nil
}