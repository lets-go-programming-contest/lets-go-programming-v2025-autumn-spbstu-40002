package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrCantDecorate           = errors.New("can't be decorated")
	ErrContextDoneInDecorator = errors.New("context done in decorator")
	ErrContextDoneInSeparator = errors.New("context done in separator")
)

const (
	StrNoDecorator = "no decorator"
	StrNoMult      = "no multiplexer"
	StrDecorated   = "decorated: "
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, StrNoDecorator) {
				return ErrCantDecorate
			}

			if !strings.HasPrefix(data, StrDecorated) {
				data = StrDecorated + data
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return errors.Join(ErrContextDoneInDecorator, ctx.Err())
			}
		case <-ctx.Done():
			return errors.Join(ErrContextDoneInDecorator, ctx.Err())
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	index := 0

	for {
		select {
		case data, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[index] <- data:
			case <-ctx.Done():
				return errors.Join(ErrContextDoneInSeparator, ctx.Err())
			}

			index = (index + 1) % len(outputs)
		case <-ctx.Done():
			return errors.Join(ErrContextDoneInSeparator, ctx.Err())
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	waitGroup := &sync.WaitGroup{}
	done := make(chan struct{})

	for _, inputChannel := range inputs {
		waitGroup.Add(1)

		go func(inputChan chan string) {
			defer waitGroup.Done()

			for {
				select {
				case data, ok := <-inputChan:
					if !ok {
						return
					}

					if strings.Contains(data, StrNoMult) {
						continue
					}

					select {
					case output <- data:
					case <-ctx.Done():
						return
					case <-done:
						return
					}
				case <-ctx.Done():
					return
				case <-done:
					return
				}
			}
		}(inputChannel)
	}

	<-ctx.Done()
	close(done)
	waitGroup.Wait()

	return nil
}
