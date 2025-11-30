package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCannotDecorate = errors.New("can't be decorated")

func PrefixDecoratorFunc(
	ctx context.Context,
	inputChan chan string,
	outputChan chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case value, ok := <-inputChan:
			if !ok {
				return nil
			}

			if strings.Contains(value, "no decorator") {
				return ErrCannotDecorate
			}

			if !strings.HasPrefix(value, "decorated: ") {
				value = "decorated: " + value
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case outputChan <- value:
			}
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	inputChan chan string,
	outputChans []chan string,
) error {
	if len(outputChans) == 0 {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case _, ok := <-inputChan:
				if !ok {
					return nil
				}
			}
		}
	}

	idx := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case value, ok := <-inputChan:
			if !ok {
				return nil
			}

			outputChan := outputChans[idx%len(outputChans)]
			idx++

			select {
			case <-ctx.Done():
				return ctx.Err()
			case outputChan <- value:
			}
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputChans []chan string,
	outputChan chan string,
) error {
	var waitGroup sync.WaitGroup

	for _, channel := range inputChans {
		waitGroup.Add(1)
		ch := channel
		go func() {
			defer waitGroup.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case value, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(value, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case outputChan <- value:
					}
				}
			}
		}()
	}

	waitGroup.Wait()
	return nil
}
