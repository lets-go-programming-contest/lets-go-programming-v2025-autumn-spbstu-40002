package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrCantDecorate  = errors.New("can't be decorated")
	ErrNoMultiplexer = errors.New("no multiplexer")
	ErrNoSeparator   = errors.New("no separator")
)

func PrefixDecoratorFunc(ctx context.Context, in chan string, out chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case v, ok := <-in:
			if !ok {
				return nil
			}
			if strings.Contains(v, "no decorator") {
				return ErrCantDecorate
			}
			select {
			case <-ctx.Done():
				return nil
			case out <- "decorated: " + v:
			}
		}
	}
}

func consumeOne(
	ctx context.Context,
	in chan string,
	out chan string,
	errc chan<- error,
	stop <-chan struct{},
) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-stop:
			return
		case v, ok := <-in:
			if !ok {
				return
			}
			if strings.Contains(v, "no multiplexer") {
				select {
				case errc <- ErrNoMultiplexer:
				default:
				}
				return
			}
			select {
			case <-ctx.Done():
				return
			case <-stop:
				return
			case out <- v:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, out chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup

	errc := make(chan error, 1)
	stop := make(chan struct{})

	wg.Add(len(inputs))

	for _, ch := range inputs {
		in := ch

		go func() {
			defer wg.Done()
			consumeOne(ctx, in, out, errc, stop)
		}()
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		close(stop)
		<-done
		return nil
	case <-done:
		select {
		case err := <-errc:
			return err
		default:
			return nil
		}
	}
}

func SeparatorFunc(ctx context.Context, in chan string, outs []chan string) error {
	if len(outs) == 0 {
		return nil
	}

	idx := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case v, ok := <-in:
			if !ok {
				return nil
			}
			if strings.Contains(v, "no separator") {
				return ErrNoSeparator
			}

			target := outs[idx]

			select {
			case <-ctx.Done():
				return nil
			case target <- v:
			}

			idx++
			if idx == len(outs) {
				idx = 0
			}
		}
	}
}
