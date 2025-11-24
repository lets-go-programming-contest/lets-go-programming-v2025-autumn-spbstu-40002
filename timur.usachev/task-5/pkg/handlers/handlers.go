package handlers

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case val, ok := <-input:
			if !ok {
				close(output)
				return nil
			}
			if strings.Contains(val, "no decorator") {
				return fmt.Errorf("can't be decorated")
			}
			if !strings.HasPrefix(val, "decorated: ") {
				val = "decorated: " + val
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- val:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case _, ok := <-input:
				if !ok {
					return nil
				}
			}
		}
	}
	counter := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case val, ok := <-input:
			if !ok {
				for _, out := range outputs {
					if out != nil {
						close(out)
					}
				}
				return nil
			}
			out := outputs[counter%len(outputs)]
			counter++
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- val:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		close(output)
		return nil
	}
	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, ch := range inputs {
		ch := ch
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-ch:
					if !ok {
						return
					}
					if strings.Contains(val, "no multiplexer") {
						continue
					}
					select {
					case <-ctx.Done():
						return
					case output <- val:
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
		return ctx.Err()
	case <-done:
		close(output)
		return nil
	}
}
