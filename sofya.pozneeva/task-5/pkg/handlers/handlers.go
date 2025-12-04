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
	defer close(output)

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
	defer (func() {
		for _, ch := range outputs {
			close(ch)
		}
	})()

	var (
		i   = 0
		cntOut = len(outputs)
	)
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
            
        case value, ok := <-input:
            if !ok {
                return nil
            }
            
            select {
            case outputs[i] <- value:
                i = (i + 1) % cntOut
                
            case <-ctx.Done():
                return ctx.Err()
            }
        }
    }
}
