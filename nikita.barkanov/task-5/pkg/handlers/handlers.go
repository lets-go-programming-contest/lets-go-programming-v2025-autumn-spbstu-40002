package handlers

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(data, "no decorator") {
				return fmt.Errorf("can’t be decorated")
			}
			if !strings.HasPrefix(data, "decorated: ") {
				data = "decorated: " + data
			}
			output <- data
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, out := range outputs {
			close(out)
		}
	}()

	i := 0
	numOutputs := len(outputs)
	if numOutputs == 0 {
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}
			outputs[i%numOutputs] <- data
			i++
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, inCh := range inputs {
		go func(ch <-chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-ch:
					if !ok {
						return
					}
					if strings.Contains(data, "no multiplexer") {
						continue
					}
					select {
					case output <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(inCh)
	}

	wg.Wait()
	return ctx.Err()
}
