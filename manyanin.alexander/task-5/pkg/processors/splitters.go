package processors

import (
	"context"
	"fmt"
)

var noOutputChannels = fmt.Errorf("output channels missing")

func RoundRobin(
	ctx context.Context,
	inCh chan string,
	outChs []chan string,
) error {
	if len(outChs) == 0 {
		return noOutputChannels
	}

	defer func() {
		for _, ch := range outChs {
			close(ch)
		}
	}()

	counter := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-inCh:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case outChs[counter] <- data:
			}

			counter = (counter + 1) % len(outChs)
		}
	}
}

func ByLength(
	ctx context.Context,
	inCh chan string,
	outChs []chan string,
) error {
	if len(outChs) < 2 {
		return noOutputChannels
	}

	defer func() {
		for _, ch := range outChs {
			close(ch)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-inCh:
			if !ok {
				return nil
			}

			idx := 0
			if len(data) > 5 {
				idx = 1
			}

			if idx < len(outChs) {
				select {
				case <-ctx.Done():
					return nil
				case outChs[idx] <- data:
				}
			}
		}
	}
}
