package main

import (
	"context"
	"log"
	"time"

	"github.com/t1wt/task-5/pkg/conveyer"
	"github.com/t1wt/task-5/pkg/handlers"
)

const (
	bufferSize     = 4
	recvTimeoutSec = 2
	repeatCount    = 3
)

func runConveyer(ctx context.Context, cvr *conveyer.Conveyer) <-chan error {
	errCh := make(chan error, 1)

	go func() {
		errCh <- cvr.Run(ctx)
	}()

	return errCh
}

func sendInputs(cvr *conveyer.Conveyer, inputs []string) {
	for _, item := range inputs {
		if err := cvr.Send("input", item); err != nil {
			log.Println(err)
		}
	}
}

func receiveOutputs(cvr *conveyer.Conveyer) {
	for range make([]struct{}, repeatCount) {
		resCh := make(chan struct{})

		var (
			val string
			err error
		)

		go func() {
			defer close(resCh)

			val, err = cvr.Recv("out")
		}()

		select {
		case <-resCh:
			if err != nil {
				log.Println(err)

				continue
			}

			log.Println("recv:", val)

		case <-time.After(recvTimeoutSec * time.Second):
			log.Println("timeout")
		}
	}
}

func main() {
	cvr := conveyer.New(bufferSize)

	cvr.RegisterDecorator(handlers.PrefixDecoratorFunc, "input", "decorated")
	cvr.RegisterSeparator(handlers.SeparatorFunc, "decorated", []string{"s1", "s2"})
	cvr.RegisterMultiplexer(handlers.MultiplexerFunc, []string{"s1", "s2"}, "out")

	ctx := context.Background()
	ctxRun, cancelRun := context.WithCancel(ctx)

	defer cancelRun()

	runErrCh := runConveyer(ctxRun, cvr)

	inputs := []string{
		"hello",
		"no multiplexer should be filtered",
		"world",
		"no decorator please",
		"gopher",
	}

	sendInputs(cvr, inputs)
	receiveOutputs(cvr)

	cancelRun()

	if err := <-runErrCh; err != nil {
		log.Println("conveyer finished with error:", err)

		return
	}

	log.Println("conveyer finished")
}
