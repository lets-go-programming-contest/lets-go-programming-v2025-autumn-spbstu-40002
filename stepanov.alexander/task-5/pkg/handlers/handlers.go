package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrNoDecorator  = errors.New("can't be decorated")
	ErrEmptyOutputs = errors.New("empty outputs")
)

const (
	decoratedPrefix = "decorated: "
	noDecoratorMark = "no decorator"
	noMuxMark       = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, inputChannel, outputChannel chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-inputChannel:
			if !ok {
				return nil
			}

			if strings.Contains(data, noDecoratorMark) {
				return ErrNoDecorator
			}

			if !strings.HasPrefix(data, decoratedPrefix) {
				data = decoratedPrefix + data
			}

			select {
			case outputChannel <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputChannels []chan string, outputChannel chan string) error {
	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputChannels))

	for _, inputChannel := range inputChannels {
		currentInputChannel := inputChannel

		go func() {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case data, ok := <-currentInputChannel:
					if !ok {
						return
					}

					if strings.Contains(data, noMuxMark) {
						continue
					}

					select {
					case outputChannel <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	waitGroup.Wait()

	return nil
}

func SeparatorFunc(ctx context.Context, inputChannel chan string, outputChannels []chan string) error {
	if len(outputChannels) == 0 {
		return ErrEmptyOutputs
	}

	outputIndex := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-inputChannel:
			if !ok {
				return nil
			}

			select {
			case outputChannels[outputIndex] <- data:
			case <-ctx.Done():
				return nil
			}

			outputIndex = (outputIndex + 1) % len(outputChannels)
		}
	}
}
