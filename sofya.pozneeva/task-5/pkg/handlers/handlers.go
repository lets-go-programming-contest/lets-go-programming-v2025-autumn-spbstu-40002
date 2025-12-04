package handlers

import (
	"context"
	"errors"
	"strings"
)

var (
	errStringNotDecorated = errors.New("can’t be decorated")
	errStringDecorated    = errors.New("string was decorated earlier")
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

		case outputString, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(inputString, noDecorated) {
				return errStringNotDecorated
			}

			if !strings.Contains(inputString, stringForDecorate) {
				outputString = stringForDecorate + outputString
			}

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
					return nil
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
		i      = 0
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
            case outputs[i] <- value:
                i = (i + 1) % cntOut
                
            case <-ctx.Done():
                return nil
            }
        }
    }
}
