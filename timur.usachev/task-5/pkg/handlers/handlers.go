package handlers

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case val, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(val, "no decorator") {
				return fmt.Errorf("can't be decorated")
			}
			if !strings.HasPrefix(val, "decorated: ") {
				val = "decorated: " + val
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- val:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
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
	counter := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case val, ok := <-input:
			if !ok {
				return nil
			}
			out := outputs[counter%len(outputs)]
			counter++
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- val:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	cases := make([]reflect.SelectCase, 0, len(inputs))
	for _, ch := range inputs {
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)})
	}
	closedCount := 0
	for {
		if closedCount == len(cases) {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		chosen, recv, ok := reflect.Select(cases)
		if !ok {
			cases[chosen].Chan = reflect.ValueOf(nil)
			closedCount++
			continue
		}
		val := recv.Interface().(string)
		if strings.Contains(val, "no multiplexer") {
			continue
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case output <- val:
		}
	}
}
