package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/megurumacabre/task-5/pkg/conveyer"
)

type DecoratorFunc = conveyer.DecoratorFunc
type MultiplexerFunc = conveyer.MultiplexerFunc
type SeparatorFunc = conveyer.SeparatorFunc

var (
	ErrDecorator   = errors.New("decorator failed")
	ErrMultiplexer = errors.New("multiplexer failed")
	ErrSeparator   = errors.New("separator failed")
)

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
				return fmt.Errorf("no decorator: %w", ErrDecorator)
			}

			outVal := "decorated: " + value

			select {
			case <-ctx.Done():
				return nil
			case output <- outVal:
			}
		}
	}
}

func consumeOne(ctx context.Context, input chan string, output chan string, errOnce chan<- error, stop <-chan struct{}) {
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
				case errOnce <- fmt.Errorf("no multiplexer: %w", ErrMultiplexer):
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

func MergeMultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	var waitGroup sync.WaitGroup

	errOnce := make(chan error, 1)
	stopAll := make(chan struct{})

	for _, ch := range inputs {
		inCh := ch

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()
			consumeOne(ctx, inCh, output, errOnce, stopAll)
		}()
	}

	done := make(chan struct{})

	go func() {
		waitGroup.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		close(stopAll)
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

func RoundRobinSeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
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
				return fmt.Errorf("no separator: %w", ErrSeparator)
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
