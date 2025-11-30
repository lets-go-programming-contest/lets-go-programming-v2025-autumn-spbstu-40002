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

const (
	sleepDuration   = 10 * time.Millisecond
	timeoutDuration = 100 * time.Millisecond
)

func readBufferSize() (int, error) {
	var bufferSize int
	_, err := fmt.Scan(&bufferSize)

	if err != nil {
		return 0, fmt.Errorf("failed to read buffer size: %w", err)
	}

	return bufferSize, nil
}

func readNumHandlers() (int, error) {
	var numHandlers int
	_, err := fmt.Scan(&numHandlers)

	if err != nil {
		return 0, fmt.Errorf("failed to read number of handlers: %w", err)
	}

	return numHandlers, nil
}

func registerDecorator(conv *conveyer.Conveyer) error {
	var input, output string
	_, err := fmt.Scan(&input, &output)

	if err != nil {
		return fmt.Errorf("failed to read decorator channels: %w", err)
	}

	conv.RegisterDecorator(handlers.PrefixDecoratorFunc, input, output)

	return nil
}

func registerMultiplexer(conv *conveyer.Conveyer) error {
	var numInputs int
	_, err := fmt.Scan(&numInputs)

	if err != nil {
		return fmt.Errorf("failed to read multiplexer inputs count: %w", err)
	}

	inputs := make([]string, numInputs)
	for i := range numInputs {
		_, err = fmt.Scan(&inputs[i])
		if err != nil {
			return fmt.Errorf("failed to read multiplexer input: %w", err)
		}
	}

	var output string
	_, err = fmt.Scan(&output)
	if err != nil {
		return fmt.Errorf("failed to read multiplexer output: %w", err)
	}

	conv.RegisterMultiplexer(handlers.MultiplexerFunc, inputs, output)

	return nil
}

func registerSeparator(conv *conveyer.Conveyer) error {
	var input string

	var numOutputs int
	_, err := fmt.Scan(&input, &numOutputs)

	if err != nil {
		return fmt.Errorf("failed to read separator channels: %w", err)
	}

	outputs := make([]string, numOutputs)
	for i := range numOutputs {
		_, err = fmt.Scan(&outputs[i])
		if err != nil {
			return fmt.Errorf("failed to read separator output: %w", err)
		}
	}

	conv.RegisterSeparator(handlers.SeparatorFunc, input, outputs)

	return nil
}

func registerHandlers(conv *conveyer.Conveyer, numHandlers int) error {
	for i := range numHandlers {
		_ = i // intrange: use integer range when possible

		var handlerType string
		_, err := fmt.Scan(&handlerType)

		if err != nil {
			return fmt.Errorf("failed to read handler type: %w", err)
		}

		switch handlerType {
		case "decorator":
			if err := registerDecorator(conv); err != nil {
				return err
			}

		case "multiplexer":
			if err := registerMultiplexer(conv); err != nil {
				return err
			}

		case "separator":
			if err := registerSeparator(conv); err != nil {
				return err
			}
		}
	}

	return nil
}

func readChannelNames() (string, string, error) {
	var inputChan string
	_, err := fmt.Scan(&inputChan)

	if err != nil {
		return "", "", fmt.Errorf("failed to read input channel: %w", err)
	}

	var outputChan string
	_, err = fmt.Scan(&outputChan)

	if err != nil {
		return "", "", fmt.Errorf("failed to read output channel: %w", err)
	}

	return inputChan, outputChan, nil
}

func sendData(conv *conveyer.Conveyer, inputChan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data := scanner.Text()
		if sendErr := conv.Send(inputChan, data); sendErr != nil {
			break
		}
	}
}

func receiveAndOutput(conv *conveyer.Conveyer, outputChan string, sigChan chan os.Signal, errChan chan error) bool {
	select {
	case <-sigChan:
		return true
	case err := <-errChan:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

		return true
	default:
		result, recvErr := conv.Recv(outputChan)
		if recvErr != nil {
			time.Sleep(sleepDuration)

			return false
		}

		if result == "undefined" {
			time.Sleep(sleepDuration)

			return false
		}

		_, _ = fmt.Fprintf(os.Stdout, "%s\n", result)

		return false
	}
}

func main() {
	bufferSize, err := readBufferSize()
	if err != nil {
		return
	}

	conv := conveyer.New(bufferSize)

	numHandlers, err := readNumHandlers()
	if err != nil {
		return
	}

	err = registerHandlers(conv, numHandlers)
	if err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		errChan <- conv.Run(ctx)
	}()

	inputChan, outputChan, err := readChannelNames()
	if err != nil {
		cancel()
		<-errChan

		return
	}

	go sendData(conv, inputChan)

	done := false
	for !done {
		done = receiveAndOutput(conv, outputChan, sigChan, errChan)
	}

	select {
	case err := <-errChan:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

	case <-time.After(timeoutDuration):
	}
}
