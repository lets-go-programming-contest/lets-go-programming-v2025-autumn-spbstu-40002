package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrCannotDecorate = errors.New("can't be decorated")
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case v, ok := <-input:
			if !ok {
				close(output)
				return nil
			}

			if strings.Contains(v, "no decorator") {
				return ErrCannotDecorate
			}

			if !strings.HasPrefix(v, "decorated: ") {
				v = "decorated: " + v
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- v:
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

	idx := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case v, ok := <-input:
			if !ok {
				// ВАЖНО: закрыть все выходные каналы
				for _, ch := range outputs {
					close(ch)
				}
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

	for _, ch := range inputs {
		wg.Add(1)

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

	go func() {
		wg.Wait()
		close(output)
	}()

	<-ctx.Done()
	return ctx.Err()
}
