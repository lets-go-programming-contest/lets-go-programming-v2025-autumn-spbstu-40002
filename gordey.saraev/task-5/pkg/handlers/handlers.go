package conveyer

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrCantBeDecorated = errors.New("can't be decorated")
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	const prefix = "decorated: "

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrCantBeDecorated
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
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	counter := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			idx := counter % len(outputs)
			counter++

			select {
			case <-ctx.Done():
				return nil
			case outputs[idx] <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	done := make(chan struct{})
	defer close(done)

	type dataWithIndex struct {
		data  string
		index int
		ok    bool
	}

	merged := make(chan dataWithIndex)

	for i, input := range inputs {
		go func(idx int, ch chan string) {
			for {
				select {
				case <-done:
					return
				case <-ctx.Done():
					return
				case data, ok := <-ch:
					select {
					case <-done:
						return
					case <-ctx.Done():
						return
					case merged <- dataWithIndex{data, idx, ok}:
						if !ok {
							return
						}
					}
				}
			}
		}(i, input)
	}

	activeInputs := len(inputs)
	for {
		select {
		case <-ctx.Done():
			return nil
		case dwi := <-merged:
			if !dwi.ok {
				activeInputs--
				if activeInputs == 0 {
					return nil
				}
				continue
			}

			if strings.Contains(dwi.data, "no multiplexer") {
				continue
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- dwi.data:
			}
		}
	}
}
