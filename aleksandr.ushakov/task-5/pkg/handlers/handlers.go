package handlers

import (
	"context"
	"strings"
)

const (
	noDecoratorMsg   = "no decorator"
	msgForDecorator  = "decorated: "
	noMultiplexerMsg = "no multiplexer"
)

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, noDecoratorMsg) {
				return ErrNoDecorate
			}

			if !strings.HasPrefix(data, msgForDecorator) {
				data = msgForDecorator + data
			}
			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	defer close(output)

	for _, inputCh := range inputs {
		go func(channel <-chan string) {
			for {
				select {
				case <-ctx.Done():
					return
				case data, chanOpen := <-channel:
					if !chanOpen {
						return
					}

					if presenceTag := strings.Contains(data, noMultiplexerMsg); !presenceTag {
						select {
						case <-ctx.Done():
							return
						case output <- data:
						}
					}
				}
			}
		}(inputCh)
	}

	<-ctx.Done()

	return nil
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	defer func() {
		for _, out := range outputs {
			close(out)
		}
	}()

	if len(outputs) == 0 {
		return nil
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[index] <- data:
				index = (index + 1) % len(outputs)
			case <-ctx.Done():
				return nil
			}
		}
	}
}
