package handlers

import (
	"context"
	"errors"
	"strings"
)

var errStringNotDecorated = errors.New("can't be decorated")

const (
	noMultiplexerData = "no multiplexer"
	noDecorated       = "no decorator"
	stringForDecorate = "decorated: "
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
		case outputString, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(outputString, noDecorated) {
				return errStringNotDecorated
			}

			flag := strings.Contains(outputString, stringForDecorate)
			if !flag {
				outputString = stringForDecorate + outputString
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- outputString:
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

	if err := ctx.Err(); err != nil {
		return nil
	}

	for _, channel := range inputs {
		go func(ch <-chan string) {
			for {
				select {
				case <-ctx.Done():
					return
				case value, ok := <-ch:
					if !ok {
						return
					}

					if !strings.Contains(value, noMultiplexerData) {
						select {
						case <-ctx.Done():
							return
						case output <- value:
						}
					}
				}
			}
		}(channel)
	}

	<-ctx.Done()

	return nil
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	defer (func() {
		for _, ch := range outputs {
			close(ch)
		}
	})()

	var (
		cnt    = 0
		cntOut = len(outputs)
	)

	for {
		select {
		case <-ctx.Done():
			return nil
		case value, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case outputs[cnt] <- value:
				cnt = (cnt + 1) % cntOut
			}
		}
	}
}
