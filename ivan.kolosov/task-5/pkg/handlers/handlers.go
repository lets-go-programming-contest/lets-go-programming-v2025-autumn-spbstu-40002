package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

const (
	noDecorator   = "no decorator"
	decoratorText = "decorated: "
	noMultiplexer = "no multiplexer"
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
		case str, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(str, noDecorator) {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(str, decoratorText) {
				str = decoratorText + str
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- str:
			}
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	counter := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case str, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case outputs[counter] <- str:
			}

			counter++
			if counter >= len(outputs) {
				counter = 0
			}
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	var wgr sync.WaitGroup

	wgr.Add(len(inputs))

	for _, curChannel := range inputs {
		go func(channel <-chan string) {
			defer wgr.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case str, ok := <-channel:
					if !ok {
						return
					}

					if strings.Contains(str, noMultiplexer) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- str:
					}
				}
			}
		}(curChannel)
	}

	wgr.Wait()

	return nil
}
