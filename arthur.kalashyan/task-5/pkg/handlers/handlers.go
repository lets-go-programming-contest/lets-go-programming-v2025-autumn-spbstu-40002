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
		case val, ok := <-inputChan:
			if !ok {
				return nil
			}

			if strings.Contains(val, "no decorator") {
				return ErrCannotDecorate
			}

			if !strings.HasPrefix(val, "decorated: ") {
				val = "decorated: " + val
			}

			select {
			case <-ctx.Done():
				return errors.Join(ctx.Err())
			case outputChan <- val:
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

	idx := 0
	for {
		select {
		case <-ctx.Done():
			return errors.Join(ctx.Err())
		case val, ok := <-inputChan:
			if !ok {
				return nil
			}

			out := outputChans[idx%len(outputChans)]
			idx++

			select {
			case <-ctx.Done():
				return errors.Join(ctx.Err())
			case out <- val:
			}
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputChans []chan string,
	outputChan chan string,
) error {
	var wg sync.WaitGroup

	for _, ch := range inputChans {
		chCopy := ch
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-chCopy:
					if !ok {
						return
					}

					if strings.Contains(val, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case outputChan <- val:
					}
				}
			}
		}()
	}

	wg.Wait()
	return nil
}
