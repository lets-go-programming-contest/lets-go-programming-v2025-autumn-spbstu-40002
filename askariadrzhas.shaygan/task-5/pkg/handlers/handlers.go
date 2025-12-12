package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrDecorationFailed = errors.New("decoration not possible")
	ErrNoSources        = errors.New("no sources provided")
	ErrNoDestinations   = errors.New("no destinations provided")
)

const (
	SkipDecorator  = "no decorator"
	SkipMerger     = "no multiplexer"
	DecoratedLabel = "decorated: "
)

func PrefixProcessor(ctx context.Context, inCh chan string, outCh chan string) error {
	for {
		select {
		case data, active := <-inCh:
			if !active {
				return nil
			}

			if strings.Contains(data, SkipDecorator) {
				return ErrDecorationFailed
			}

			result := data
			if !strings.HasPrefix(data, DecoratedLabel) {
				result = DecoratedLabel + data
			}

			select {
			case <-ctx.Done():
				return nil
			case outCh <- result:
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func SplitProcessor(ctx context.Context, inCh chan string, outChs []chan string) error {
	if len(outChs) == 0 {
		return ErrNoDestinations
	}

	idx := 0

	for {
		select {
		case data, active := <-inCh:
			if !active {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case outChs[idx] <- data:
			}

			idx = (idx + 1) % len(outChs)
		case <-ctx.Done():
			return nil
		}
	}
}

func MergeProcessor(ctx context.Context, inChs []chan string, outCh chan string) error {
	if len(inChs) == 0 {
		return ErrNoSources
	}

	var wg sync.WaitGroup
	wg.Add(len(inChs))

	for _, ch := range inChs {
		go func(input chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, active := <-input:
					if !active {
						return
					}

					if !strings.Contains(data, SkipMerger) {
						select {
						case <-ctx.Done():
							return
						case outCh <- data:
						}
					}
				}
			}
		}(ch)
	}

	wg.Wait()
	return nil
}
