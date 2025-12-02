package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrNoDecorator = errors.New("can't be decorated")

const (
	noDecoratorData        = "no decorator"
	textForDecoratorString = "decorated: "
	noMultiplexerData      = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil

		case line, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(line, noDecoratorData) {
				return ErrNoDecorator
			}

			if !strings.HasPrefix(line, textForDecoratorString) {
				line = textForDecoratorString + line
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- line:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	for _, inputCh := range inputs {
		inCh := inputCh

		go func(chLocal chan string) {
			for {
				select {
				case <-ctx.Done():
					return

				case line, ok := <-chLocal:
					if !ok {
						return
					}

					if strings.Contains(line, noMultiplexerData) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- line:
					}
				}
			}
		}(inCh)
	}

	<-ctx.Done()

	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, outCh := range outputs {
			close(outCh)
		}
	}()

	countOut := len(outputs)
	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case line, ok := <-input:
			if !ok {
				return nil
			}

			if countOut == 0 {
				continue
			}

			select {
			case <-ctx.Done():
				return nil
			case outputs[index%countOut] <- line:
				index++
			}
		}
	}
}
