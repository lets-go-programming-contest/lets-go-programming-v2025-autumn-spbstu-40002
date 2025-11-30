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
			return errors.Join(ctx.Err())
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
				return errors.Join(ctx.Err())
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
				return errors.Join(ctx.Err())
			case _, ok := <-inputChan:
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
			return errors.Join(ctx.Err())
		case value, ok := <-inputChan:
			if !ok {
				return nil
			}

			outputChan := outputChans[index%len(outputChans)]
			index++

			select {
			case <-ctx.Done():
				return errors.Join(ctx.Err())
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
		channelCopy := channel

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case value, ok := <-channelCopy:
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
