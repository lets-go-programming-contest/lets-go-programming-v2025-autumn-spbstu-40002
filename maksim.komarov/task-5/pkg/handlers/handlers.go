package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

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
		case v, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(v, "no decorator") {
				return fmt.Errorf("no decorator: %w", ErrDecorator)
			}

			outVal := "decorated: " + v

			select {
			case <-ctx.Done():
				return nil
			case output <- outVal:
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
		case v, ok := <-input:
			if !ok {
				return
			}

			if strings.Contains(v, "no multiplexer") {
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
			case output <- v:
			}
		}
	}
}

func MergeMultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup

	errOnce := make(chan error, 1)
	stopAll := make(chan struct{})

	for _, ch := range inputs {
		in := ch

		wg.Add(1)

		go func() {
			defer wg.Done()
			consumeOne(ctx, in, output, errOnce, stopAll)
		}()
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
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

	i := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case v, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(v, "no separator") {
				return fmt.Errorf("no separator: %w", ErrSeparator)
			}

			target := outputs[i]

			select {
			case <-ctx.Done():
				return nil
			case target <- v:
			}

			i++

			if i == len(outputs) {
				i = 0
			}
		}
	}
}

var MultiplexerFunc = MergeMultiplexerFunc
var SeparatorFunc = RoundRobinSeparatorFunc
