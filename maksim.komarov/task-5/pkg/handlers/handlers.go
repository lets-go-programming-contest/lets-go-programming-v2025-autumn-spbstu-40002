package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCantDecorate = errors.New("can't be decorated")

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(msg, "no decorator") {
				return ErrCantDecorate
			}

			decorated := "decorated: " + msg

			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- decorated:
			}
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	var wg sync.WaitGroup

	wg.Add(len(inputs))

	for _, ch := range inputs {
		in := ch

		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-in:
					if !ok {
						return
					}

					if strings.Contains(msg, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- msg:
					}
				}
			}
		}()
	}

	wg.Wait()

	return nil
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	if len(outputs) == 0 {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case _, ok := <-input:
				if !ok {
					return nil
				}
			}
		}
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-input:
			if !ok {
				return nil
			}

			dst := outputs[index%len(outputs)]
			index++

			select {
			case <-ctx.Done():
				return ctx.Err()
			case dst <- msg:
			}
		}
	}
}
