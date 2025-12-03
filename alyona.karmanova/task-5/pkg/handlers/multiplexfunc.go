package handlers

import (
	"context"
	"fmt"
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

	var waitg sync.WaitGroup

	waitg.Add(len(inputs))

	for _, inputit := range inputs {
		go func() {
			defer waitg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case value, ok := <-inputit:
					if !ok {
						return
					}

					if strings.Contains(value, "no multiplexer") {
						continue
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

	waitg.Wait()

	return fmt.Errorf("multiplexer finished: %w", ctx.Err())

}
