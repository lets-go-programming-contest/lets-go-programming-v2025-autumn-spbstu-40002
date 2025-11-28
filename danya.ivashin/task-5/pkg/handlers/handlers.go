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
	noDecorator     = "no decorator"
	noMultiplexer   = "no multiplexer"
	decoratedPrefix = "decorated: "
)

func orDone[T any](done <-chan struct{}, channel <-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			default:
			}

			select {
			case value, ok := <-channel:
				if !ok {
					return
				}

				out <- value

			default:
			}
		}
	}()

	return out
}

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	for str := range orDone(ctx.Done(), input) {
		if strings.Contains(str, noDecorator) {
			return ErrNoDecorator
		}

		if !strings.HasPrefix(str, decoratedPrefix) {
			str = decoratedPrefix + str
		}

		select {
		case output <- str:

		case <-ctx.Done():
			return nil
		}
	}

	return nil
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	if len(outputs) == 0 {
		return ErrEmptyOutputs
	}

	index := 0

	for str := range orDone(ctx.Done(), input) {
		select {
		case outputs[index] <- str:

		case <-ctx.Done():
			return nil
		}

		index = (index + 1) % len(outputs)
	}

	return nil
}

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	waitgr := sync.WaitGroup{}

	waitgr.Add(len(inputs))

	for _, ch := range inputs {
		go func() {
			defer waitgr.Done()

			for str := range orDone(ctx.Done(), ch) {
				if strings.Contains(str, noMultiplexer) {
					continue
				}

				select {
				case output <- str:

				case <-ctx.Done():
					return
				}
			}
		}()
	}

	waitgr.Wait()

	return nil
}
