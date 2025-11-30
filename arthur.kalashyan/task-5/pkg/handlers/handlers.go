package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case v, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(v, "no decorator") {
				return errors.New("can't be decorated: " + v)
			}
			var out string
			if strings.HasPrefix(v, "decorated: ") {
				out = v
			} else {
				out = "decorated: " + v
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- out:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		for range input {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}
		return nil
	}
	idx := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case v, ok := <-input:
			if !ok {
				return nil
			}
			outCh := outputs[idx%len(outputs)]
			idx++
			select {
			case <-ctx.Done():
				return ctx.Err()
			case outCh <- v:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)
	for _, ch := range inputs {
		wg.Add(1)
		c := ch
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					if strings.Contains(v, "no multiplexer") {
						continue
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
	go func() {
		wg.Wait()
		close(errCh)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case _, ok := <-errCh:
		if !ok {
			return nil
		}
		return nil
	}
}
