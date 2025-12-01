package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrCantBeDecotated = errors.New("can't be decorated")

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

			if strings.Contains(str, "no decorator") {
				return fmt.Errorf("string \"%s\" %w", str, ErrCantBeDecotated)
			}

			if !strings.HasPrefix(str, "decorated:") {
				str = "decorated: " + str
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
	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputs))

	for _, chanel := range inputs {
		go func(chanel chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case str, ok := <-chanel:
					if !ok {
						return
					}

					if strings.Contains(str, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- str:
					}
				}
			}
		}(chanel)
	}

	waitGroup.Wait()

	return nil
}
