package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/megurumacabre/task-5/pkg/conveyer"
)

var (
	ErrCantDecorate  = errors.New("can't be decorated")
	ErrNoMultiplexer = errors.New("no multiplexer")
	ErrNoSeparator   = errors.New("no separator")
)

type DecoratorFunc = conveyer.DecoratorFunc
type MultiplexerFunc = conveyer.MultiplexerFunc
type SeparatorFunc = conveyer.SeparatorFunc

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case value, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(value, "no decorator") {
				return ErrCantDecorate
			}

			out := "decorated: " + value

			select {
			case <-ctx.Done():
				return nil
			case output <- out:
			}
		}
	}
}

func consumeOne(
	ctx context.Context,
	input chan string,
	output chan string,
	errOnce chan<- error,
	stop <-chan struct{},
) {
	for {
		select {
		case <-ctx.Done():
			return

		case <-stop:
			return

		case value, ok := <-input:
			if !ok {
				return
			}

			if strings.Contains(value, "no multiplexer") {
				select {
				case errOnce <- ErrNoMultiplexer:
				default:
				}

				return
			}

			select {
			case <-ctx.Done():
				return
			case <-stop:
				return
			case output <- value:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	var waitGroup sync.WaitGroup

	errOnce := make(chan error, 1)
	stop := make(chan struct{})

	waitGroup.Add(len(inputs))

	for _, ch := range inputs {
		inCh := ch

		go func() {
			defer waitGroup.Done()
			consumeOne(ctx, inCh, output, errOnce, stop)
		}()
	}

	done := make(chan struct{})

	go func() {
		waitGroup.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		close(stop)
		<-done

		return nil

	case <-done:
		select {
		case err := <-errOnce:
			return err
		default:
			return nil
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
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

			if strings.Contains(value, "no separator") {
				return ErrNoSeparator
			}

			target := outputs[index]

			select {
			case <-ctx.Done():
				return nil
			case target <- value:
			}

			index++

			if index == len(outputs) {
				index = 0
			}
		}
	}
}
