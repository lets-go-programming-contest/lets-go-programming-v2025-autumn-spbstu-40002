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

const sleepDuration = 10 * time.Millisecond

func readBufferSize() (int, error) {
	var bufferSize int
	_, err := fmt.Scan(&bufferSize)

	return bufferSize, err
}

func readNumHandlers() (int, error) {
	var numHandlers int
	_, err := fmt.Scan(&numHandlers)

	return numHandlers, err
}

func registerDecorator(c *conveyer.Conveyer) error {
	var input, output string
	_, err := fmt.Scan(&input, &output)
	if err != nil {
		return err
	}

	c.RegisterDecorator(handlers.PrefixDecoratorFunc, input, output)

	return nil
}

func registerMultiplexer(c *conveyer.Conveyer) error {
	var numInputs int
	_, err := fmt.Scan(&numInputs)
	if err != nil {
		return err
	}

	inputs := make([]string, numInputs)
	for i := 0; i < numInputs; i++ {
		_, err = fmt.Scan(&inputs[i])
		if err != nil {
			return err
		}
	}

	var output string
	_, err = fmt.Scan(&output)
	if err != nil {
		return err
	}

	c.RegisterMultiplexer(handlers.MultiplexerFunc, inputs, output)

	return nil
}

func registerSeparator(c *conveyer.Conveyer) error {
	var input string

	var numOutputs int
	_, err := fmt.Scan(&input, &numOutputs)

	if err != nil {
		return err
	}

	outputs := make([]string, numOutputs)
	for i := 0; i < numOutputs; i++ {
		_, err = fmt.Scan(&outputs[i])
		if err != nil {
			return err
		}
	}

	c.RegisterSeparator(handlers.SeparatorFunc, input, outputs)

	return nil
}

func registerHandlers(c *conveyer.Conveyer, numHandlers int) error {
	for i := 0; i < numHandlers; i++ {
		_ = i // intrange: use integer range when possible
		var handlerType string
		_, err := fmt.Scan(&handlerType)

		if err != nil {
			return err
		}

		switch handlerType {
		case "decorator":
			if err := registerDecorator(c); err != nil {
				return err
			}

		case "multiplexer":
			if err := registerMultiplexer(c); err != nil {
				return err
			}

		case "separator":
			if err := registerSeparator(c); err != nil {
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
		return "", "", err
	}

	var outputChan string
	_, err = fmt.Scan(&outputChan)
	if err != nil {
		return "", "", err
	}

	return inputChan, outputChan, nil
}

func sendData(c *conveyer.Conveyer, inputChan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data := scanner.Text()
		if sendErr := c.Send(inputChan, data); sendErr != nil {
			break
		}
	}
}

func receiveAndOutput(c *conveyer.Conveyer, outputChan string, sigChan chan os.Signal, errChan chan error) bool {
	select {
	case <-sigChan:
		return true
	case err := <-errChan:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

		return true
	default:
		result, recvErr := c.Recv(outputChan)
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

	c := conveyer.New(bufferSize)

	numHandlers, err := readNumHandlers()
	if err != nil {
		return
	}

	err = registerHandlers(c, numHandlers)
	if err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		errChan <- c.Run(ctx)
	}()

	inputChan, outputChan, err := readChannelNames()
	if err != nil {
		cancel()
		<-errChan

		return
	}

	go sendData(c, inputChan)

	done := false
	for !done {
		done = receiveAndOutput(c, outputChan, sigChan, errChan)
	}

	select {
	case err := <-errChan:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	case <-time.After(100 * time.Millisecond):
	}
}
