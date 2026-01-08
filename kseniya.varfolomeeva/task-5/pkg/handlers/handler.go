package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

const (
	DecorPrefix   = "decorated: "
	SkipDecorFlag = "no decorator"
	SkipMergeFlag = "no multiplexer"
)

var (
	ErrDecorFailed = errors.New("can't be decorated")
	ErrNoDest      = errors.New("outputs slice must not be empty")
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case item, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(item, SkipDecorFlag) {
				return ErrDecorFailed
			}

			if !strings.HasPrefix(item, DecorPrefix) {
				item = DecorPrefix + item
			}

			select {
			case output <- item:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return ErrNoDest
	}

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, src := range inputs {
		go func(inChan chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-inChan:
					if !ok {
						return
					}

					if strings.Contains(data, SkipMergeFlag) {
						continue
					}

					select {
					case output <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(src)
	}

	wg.Wait()
	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrNoDest
	}

	idx := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			target := outputs[idx%len(outputs)]
			idx++

			select {
			case target <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
