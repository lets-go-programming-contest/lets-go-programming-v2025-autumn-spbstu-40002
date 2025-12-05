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
			return nil
		case v, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(v, "no decorator") {
				return fmt.Errorf("can't be decorated")
			}

			if !strings.HasPrefix(v, "decorated: ") {
				v = "decorated: " + v
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- v:
			}
		}
	}

}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	i := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case v, ok := <-input:
			if !ok {
				return nil
			}

			if len(outputs) == 0 {
				continue
			}

			ch := outputs[i%len(outputs)]
			i++

			select {
			case <-ctx.Done():
				return nil
			case ch <- v:
			}
		}
	}

}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, inCh := range inputs {
		ch := inCh
		go func(c chan string) {
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
		}(ch)
	}

	wg.Wait()
	return nil

}
