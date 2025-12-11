package handlers

import (
	"context"
	"strings"
)

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	for _, inputChannel := range inputs {
		go func(channel chan string) {
			for {
				select {
				case <-ctx.Done():
					return
				case data, channelOpen := <-channel:
					if !channelOpen {
						return
					}

					if !strings.Contains(data, NoMultiplexer) {
						select {
						case <-ctx.Done():
							return
						case output <- data:
						}
					}
				}
			}
		}(inputChannel)
	}

	<-ctx.Done()

	return nil
}
