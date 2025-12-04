package handlers

import (
	"context"
	"errors"
	"strings"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)
	prefix := "decorated: "
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(data, "no decorator") {
				// ВОТ ЗДЕСЬ ИСПРАВЛЕНИЕ:
				// Нужно вернуть ошибку с текстом "can't be decorated"
				return errors.New("can't be decorated")
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

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, out := range outputs {
			close(out)
		}
	}()

	index := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}
			if len(outputs) == 0 {
				continue
			}
			select {
			case outputs[index] <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
			index = (index + 1) % len(outputs)
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	type result struct {
		data string
		ok   bool
	}

	merged := make(chan result)

	for _, in := range inputs {
		go func(in chan string) {
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-in:
					select {
					case merged <- result{data, ok}:
					case <-ctx.Done():
						return
					}
					if !ok {
						return
					}
				}
			}
		}(in)
	}

	openChannels := len(inputs)
	for openChannels > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case res := <-merged:
			if !res.ok {
				openChannels--
				continue
			}
			if strings.Contains(res.data, "no multiplexer") {
				continue
			}
			select {
			case output <- res.data:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return nil
}
