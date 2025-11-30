package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrCannotDecorate = errors.New("can't be decorated")

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("prefix decorator context error: %w", ctx.Err())
		case value, ok := <-input:
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
				return fmt.Errorf("prefix decorator context error: %w", ctx.Err())
			case output <- value:
			}
		}
	}
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
				return fmt.Errorf("separator context error: %w", ctx.Err())
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
			return fmt.Errorf("separator context error: %w", ctx.Err())
		case value, ok := <-input:
			if !ok {
				return nil
			}
			outputChan := outputs[index%len(outputs)]
			index++
			select {
			case <-ctx.Done():
				return fmt.Errorf("separator context error: %w", ctx.Err())
			case outputChan <- value:
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
	for _, inputChan := range inputs {
		waitGroup.Add(1)
		ch := inputChan
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
					case output <- value:
					}
				}
			}
		}()
	}

	waitGroup.Wait()
	return nil
}
