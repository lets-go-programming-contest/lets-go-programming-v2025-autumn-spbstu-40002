package handlers

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/sync/errgroup"
)

var ErrCantBeDecorated = errors.New("can't be decorated") // ← ASCII '

// PrefixDecoratorFunc — по ТЗ
func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	const prefix = "decorated: "
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				// Входной канал закрыт → завершаемся
				return nil
			}
			if strings.Contains(data, "no decorator") {
				return ErrCantBeDecorated
			}
			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
			}
			select {
			case output <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

// SeparatorFunc — round-robin, завершается при закрытии input
func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}
	index := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				// input закрыт → завершаемся
				return nil
			}
			target := outputs[index%len(outputs)]
			select {
			case target <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
			index++
		}
	}
}

// MultiplexerFunc — параллельно читает из всех inputs, завершается при ctx.Done()
func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Запускаем N горутин по чтению
	group, ctx := errgroup.Group{}
	for _, ch := range inputs {
		ch := ch
		group.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case data, ok := <-ch:
					if !ok {
						// этот канал закрыт — выходим из этой горутины
						return nil
					}
					if strings.Contains(data, "no multiplexer") {
						continue
					}
					select {
					case output <- data:
					case <-ctx.Done():
						return nil
					}
				}
			}
		})
	}

	// Ждём, пока все входы закроются ИЛИ отмена
	return group.Wait()
}