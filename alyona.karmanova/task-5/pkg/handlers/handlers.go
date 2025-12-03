package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrDecorator    = errors.New("can't be decorated")
	ErrOutputsEmpty = errors.New("outputs must not be empty")
)

const (
	noDecorator   = "no decorator"
	prefix        = "decorated: "
	noMultiplexer = "no multiplexer"
)

// Добавляет префикс к строкам из input и пишет в output
func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for data := range input {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if strings.Contains(data, noDecorator) {
			return ErrDecorator
		}

		if !strings.HasPrefix(data, prefix) {
			data = prefix + data
		}

		select {
		case <-ctx.Done():
			return nil
		case output <- data:
		}
	}
	return nil
}

// Разбивает входной канал на несколько выходных по кругу
func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrOutputsEmpty
	}

	index := 0
	for data := range input {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		outCh := outputs[index%len(outputs)]
		index++

		select {
		case <-ctx.Done():
			return nil
		case outCh <- data:
		}
	}
	return nil
}

// Объединяет несколько входных каналов в один выходной
func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup

	for _, in := range inputs {
		wg.Add(1)
		go func(ch chan string) {
			defer wg.Done()
			for data := range ch {
				if strings.Contains(data, noMultiplexer) {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case output <- data:
				}
			}
		}(in)
	}

	wg.Wait()
	return nil
}
