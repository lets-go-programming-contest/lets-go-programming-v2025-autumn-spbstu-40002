package handlers

import (
	"context"
	"strings"
	"sync"
)

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup
	inCh := make(chan string)
	for _, ch := range inputs {
		wg.Add(1)
		go func(c chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case inCh <- v:
					}
				}
			}
		}(ch)
	}
	go func() {
		wg.Wait()
		close(inCh)
	}()
	for {
		select {
		case <-ctx.Done():
			return nil
		case v, ok := <-inCh:
			if !ok {
				return nil
			}
			if strings.Contains(v, "no multiplexer") {
				continue
			}
			select {
			case <-ctx.Done():
				return nil
			case output <- v:
			}
		}
	}
}
