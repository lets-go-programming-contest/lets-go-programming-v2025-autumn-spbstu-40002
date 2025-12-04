package handlers

import (
	"context"
)

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, outputChannel := range outputs {
			close(outputChannel)
		}
	}()

	if len(outputs) == 0 {
		return nil
	}

	channelIndex := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, channelOpen := <-input:
			if !channelOpen {
				return nil
			}

			if len(outputs) > 0 {
				outputChannel := outputs[channelIndex]
				select {
				case <-ctx.Done():
					return nil
				case outputChannel <- data:
					channelIndex = (channelIndex + 1) % len(outputs)
				}
			}
		}
	}
}
