package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrDecorator = errors.New("decorator failed")
var ErrMultiplexer = errors.New("multiplexer failed")
var ErrSeparator = errors.New("separator failed")

func PrefixDecoratorFunc(prefix string) func(context.Context, <-chan string, chan<- string) error {
	return func(ctx context.Context, input <-chan string, output chan<- string) error {
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

				val := prefix + v

				select {
				case <-ctx.Done():
					return nil
				case output <- val:
				}
			}
		}
	}
}

func MergeMultiplexerFunc() func(context.Context, []<-chan string, chan<- string) error {
	return func(ctx context.Context, inputs []<-chan string, output chan<- string) error {
		var wg sync.WaitGroup

		errOnce := make(chan error, 1)

		for _, ch := range inputs {
			in := ch

			wg.Add(1)

			go func() {
				defer wg.Done()

				for {
					select {
					case <-ctx.Done():
						return

					case v, ok := <-in:
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
						case output <- v:
						}
					}
				}
			}()
		}

		done := make(chan struct{})

		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-ctx.Done():
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
}

func RoundRobinSeparatorFunc() func(context.Context, <-chan string, []chan<- string) error {
	return func(ctx context.Context, input <-chan string, outputs []chan<- string) error {
		if len(outputs) == 0 {
			return nil
		}

		idx := 0

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

				target := outputs[idx]

				select {
				case <-ctx.Done():
					return nil
				case target <- v:
				}

				idx++
				if idx == len(outputs) {
					idx = 0
				}
			}
		}
	}
}
