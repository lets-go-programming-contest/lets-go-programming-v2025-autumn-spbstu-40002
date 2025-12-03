package handlers

import (
	"context"
	"strings"
	"sync"
)

type Multiplexer func(
	context.Context,
	[]chan string,
	chan string,
) error

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, in := range inputs {

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

	wg.Wait()

	return ctx.Err()
}
