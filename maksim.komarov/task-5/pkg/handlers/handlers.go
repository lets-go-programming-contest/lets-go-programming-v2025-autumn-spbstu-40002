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
		case value, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(value, "no decorator") {
				return fmt.Errorf("can't be decorated")
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

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}
	var waitGroup sync.WaitGroup
	errOnce := make(chan error, 1)
	stop := make(chan struct{})

	consume := func(in chan string) {
		defer waitGroup.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-stop:
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				if strings.Contains(value, "no multiplexer") {
					select {
					case errOnce <- fmt.Errorf("no multiplexer"):
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

	waitGroup.Add(len(inputs))
	for _, ch := range inputs {
		inputChan := ch
		go consume(inputChan)
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
	idx := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case value, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(value, "no separator") {
				return fmt.Errorf("no separator")
			}
			target := outputs[idx]
			select {
			case <-ctx.Done():
				return nil
			case target <- value:
			}
			idx++
			if idx == len(outputs) {
				idx = 0
			}
		}
	}
}
