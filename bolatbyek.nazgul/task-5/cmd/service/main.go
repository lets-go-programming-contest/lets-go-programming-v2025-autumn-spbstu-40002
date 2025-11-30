package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bolatbyek/task-5/pkg/conveyer"
	"github.com/bolatbyek/task-5/pkg/handlers"
)

func main() {
	// Read buffer size
	var bufferSize int
	_, err := fmt.Scan(&bufferSize)
	if err != nil {
		return
	}

	// Create conveyer
	c := conveyer.New(bufferSize)

	// Read number of handlers
	var numHandlers int
	_, err = fmt.Scan(&numHandlers)
	if err != nil {
		return
	}

	// Register handlers
	for i := 0; i < numHandlers; i++ {
		var handlerType string
		_, err = fmt.Scan(&handlerType)
		if err != nil {
			return
		}

		switch handlerType {
		case "decorator":
			var input, output string
			_, err = fmt.Scan(&input, &output)
			if err != nil {
				return
			}
			c.RegisterDecorator(handlers.PrefixDecoratorFunc, input, output)

		case "multiplexer":
			var numInputs int
			_, err = fmt.Scan(&numInputs)
			if err != nil {
				return
			}
			inputs := make([]string, numInputs)
			for j := 0; j < numInputs; j++ {
				_, err = fmt.Scan(&inputs[j])
				if err != nil {
					return
				}
			}
			var output string
			_, err = fmt.Scan(&output)
			if err != nil {
				return
			}
			c.RegisterMultiplexer(handlers.MultiplexerFunc, inputs, output)

		case "separator":
			var input string
			var numOutputs int
			_, err = fmt.Scan(&input, &numOutputs)
			if err != nil {
				return
			}
			outputs := make([]string, numOutputs)
			for j := 0; j < numOutputs; j++ {
				_, err = fmt.Scan(&outputs[j])
				if err != nil {
					return
				}
			}
			c.RegisterSeparator(handlers.SeparatorFunc, input, outputs)
		}
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Run conveyer in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- c.Run(ctx)
	}()

	// Read input channel name
	var inputChan string
	_, err = fmt.Scan(&inputChan)
	if err != nil {
		cancel()
		<-errChan
		return
	}

	// Read output channel name
	var outputChan string
	_, err = fmt.Scan(&outputChan)
	if err != nil {
		cancel()
		<-errChan
		return
	}

	// Read and send data
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			data := scanner.Text()
			if sendErr := c.Send(inputChan, data); sendErr != nil {
				break
			}
		}
	}()

	// Receive and output data
	done := false
	for !done {
		select {
		case <-sigChan:
			cancel()
			done = true
		case err := <-errChan:
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
			done = true
		default:
			result, recvErr := c.Recv(outputChan)
			if recvErr != nil {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			if result == "undefined" {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			fmt.Println(result)
		}
	}

	// Wait for conveyer to finish
	select {
	case err := <-errChan:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	case <-time.After(100 * time.Millisecond):
		// Timeout
	}
}
