package processors

import (
	"context"
	"strings"
	"sync"
)

const skipMerge = "skip_merge"

func MergeAll(
	ctx context.Context,
	inChs []chan string,
	outCh chan string,
) error {
	if len(inChs) == 0 {
		close(outCh)
		return nil
	}

	defer close(outCh)

	var wg sync.WaitGroup
	wg.Add(len(inChs))

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	for i := range inChs {
		go func(input <-chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, open := <-input:
					if !open {
						return
					}

					if strings.Contains(data, skipMerge) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case outCh <- data:
					}
				}
			}
		}(inChs[i])
	}

	select {
	case <-ctx.Done():
		return nil
	case <-done:
		return nil
	}
}

func MergeWithFilter(
	ctx context.Context,
	inChs []chan string,
	outCh chan string,
) error {
	if len(inChs) == 0 {
		close(outCh)
		return nil
	}

	defer close(outCh)

	var wg sync.WaitGroup
	wg.Add(len(inChs))

	completion := make(chan struct{})

	go func() {
		wg.Wait()
		close(completion)
	}()

	for _, inputChan := range inChs {
		go func(ch <-chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-ch:
					if !ok {
						return
					}

					if len(data) == 0 {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case outCh <- data:
					}
				}
			}
		}(inputChan)
	}

	select {
	case <-ctx.Done():
		return nil
	case <-completion:
		return nil
	}
}
