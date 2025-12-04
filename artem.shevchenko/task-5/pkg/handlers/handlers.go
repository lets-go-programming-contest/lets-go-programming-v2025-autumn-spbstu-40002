package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrCantBeDecorated = errors.New("can’t be decorated")

// PrefixDecoratorFunc — модификатор по ТЗ
func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	const prefix = "decorated: "
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil // вход закрыт — завершаемся
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

// SeparatorFunc — round-robin по ТЗ
func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}
	idx := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}
			target := outputs[idx%len(outputs)]
			select {
			case target <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
			idx++
		}
	}
}

// MultiplexerFunc — объединяет + фильтрует по ТЗ
func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Запускаем чтение из всех входов параллельно
	for _, ch := range inputs {
		ch := ch
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-ch:
					if !ok {
						return // канал закрыт
					}
					if strings.Contains(data, "no multiplexer") {
						continue // пропускаем
					}
					select {
					case output <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	// Блокируем текущую горутину до отмены
	<-ctx.Done()
	return ctx.Err()
}