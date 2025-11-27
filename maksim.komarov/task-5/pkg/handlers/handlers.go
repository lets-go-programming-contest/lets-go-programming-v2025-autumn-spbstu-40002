package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrDecorator   = errors.New("can't be decorated")
	ErrMultiplexer = errors.New("can't be multiplexed")
	ErrSeparator   = errors.New("can't be separated")
)

func Decorator(ctx context.Context, input <-chan string, output chan<- string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case value, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(value, "no decorator") {
				return ErrDecorator
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- value:
			}
		}
	}
}

func Multiplexer(ctx context.Context, inputs []<-chan string, output chan<- string) error {
	var waitGroup sync.WaitGroup
	errChan := make(chan error, 1)

	for _, ch := range inputs {
		inputCh := ch

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case value, ok := <-inputCh:
					if !ok {
						return
					}

					if strings.Contains(value, "no multiplexer") {
						select {
						case errChan <- ErrMultiplexer:
						default:
						}

						return
					}

					select {
					case <-ctx.Done():
						return
					case output <- value:
					}
				}
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		waitGroup.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		<-done

		return nil

	case <-done:
		select {
		case err := <-errChan:
			return err
		default:
			return nil
		}
	}
}

func Separator(ctx context.Context, input <-chan string, outputs []chan<- string) error {
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
				return ErrSeparator
			}

			if len(outputs) == 0 {
				return nil
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
