package handlers

import (
	"context"
	"errors"
	"strings"
)

var (
	errStringNotDecorated = errors.New("can’t be decorated")
	errStringDecorated    = errors.New("string was decorated earlier")
	errInvalidChan        = errors.New("don't read value from chan")
)

const (
	noMultiplexerData = "no multiplexer"
	noDecorated       = "no decorator"
	stringForDecorate = "decorator:"
)

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case inputString, ok := <-input:
			if !ok {
				return errInvalidChan
			}

			if strings.Contains(inputString, stringForDecorate) {
				return errStringDecorated
			}

			if strings.Contains(inputString, noDecorated) {
				return errStringNotDecorated
			}

			decorated := stringForDecorate + inputString

			select {
			case output <- decorated:
				return nil
			case <-ctx.Done():
				return nil
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

	for {
		if err := ctx.Err(); err != nil {
			return nil
		}

		for _, channel := range inputs {
			select {
			case value, ok := <-channel:
				if !ok {
					continue
				}

				if strings.Contains(value, noMultiplexerData) {
					continue
				}
				select {
				case <-ctx.Done():
					return nil
				case output <- value:
				}
			case <-ctx.Done():
				return nil
			default:
				continue
			}
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	if ctx.Err() != nil {
		return nil
	}

	for _, channel := range outputs {
		select {
		case value, ok := <-input:
			if !ok {
				return errInvalidChan
			}
			select {
			case channel <- value:
				continue
			case <-ctx.Done():
				return nil
			}

		case <-ctx.Done():
			return nil
		}
	}

	return nil
}
