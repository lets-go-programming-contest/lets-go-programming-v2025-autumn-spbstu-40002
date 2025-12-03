package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrNoOutput    = errors.New("outputs must not be empty")
	ErrNoDecoretor = errors.New("can't be decorated")
)

const (
	noDecorator   = "no decorator"
	noMultiplexer = "no multiplexer"
	addedPrefix   = "decorated: "
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, noDecorator) {
				return ErrNoDecoretor
			}

			if !strings.HasPrefix(data, addedPrefix) {
				data = addedPrefix + data
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
	var waitGroup sync.WaitGroup

	for _, inputChanel := range inputs {
		waitGroup.Add(1)

		go func(channel chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-channel:
					if !ok {
						return
					}

					if strings.Contains(data, noMultiplexer) {
						continue
					}

					select {
					case output <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(inputChanel)
	}

	waitGroup.Wait()

	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrNoOutput
	}

	position := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case valut, ok := <-input:
			if !ok {
				return nil
			}

			outChanel := outputs[position%len(outputs)]
			position++

			select {
			case outChanel <- valut:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
