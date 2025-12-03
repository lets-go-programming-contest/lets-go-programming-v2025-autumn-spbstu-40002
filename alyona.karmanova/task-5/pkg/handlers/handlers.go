package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var errPrefixDecorator = errors.New("handlers.PrefixDecoratorFunc: can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context done: %w", ctx.Err())
		case val, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(val, "no decorator") {
				return errPrefixDecorator
			}

			if !strings.HasPrefix(val, "decorated: ") {
				val = "decorated:" + val
			}

			select {
			case <-ctx.Done():
				return fmt.Errorf("context done: %w", ctx.Err())
			case output <- val:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, in := range inputs {
		in := in
		go func() {
			defer wg.Done()
			for val := range in {
				if strings.Contains(val, "no multiplexer") {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case output <- val:
				}
			}
		}()
	}

	wg.Wait()
	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, ch := range outputs {
			close(ch)
		}
	}()

	i := 0
	for val := range input {
		out := outputs[i%len(outputs)]
		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- val:
		}
		i++
	}

	return nil
}
