package main

import (
	"context"
	"log"
	"time"

	"github.com/t1wt/task-5/pkg/conveyer"
	"github.com/t1wt/task-5/pkg/handlers"
)

const (
	chanSize       = 4
	readIterations = 3
	readTimeout    = 2 * time.Second
)

func main() {
	cvr := conveyer.New(chanSize)

	cvr.RegisterDecorator(handlers.PrefixDecoratorFunc, "input", "decorated")
	cvr.RegisterSeparator(handlers.SeparatorFunc, "decorated", []string{"s1", "s2"})
	cvr.RegisterMultiplexer(handlers.MultiplexerFunc, []string{"s1", "s2"}, "out")

	ctx := context.Background()

	ctxRun, cancelRun := context.WithCancel(ctx)
	defer cancelRun()

	runErrCh := make(chan error, 1)

	go func() {
		runErrCh <- cvr.Run(ctxRun)
	}()

	inputs := []string{
		"hello",
		"no multiplexer should be filtered",
		"world",
		"no decorator please",
		"gopher",
	}

	for _, item := range inputs {
		if err := cvr.Send("input", item); err != nil {
			log.Println(err)
		}
	}

	for range readIterations {
		resCh := make(chan struct{})

		var val string
		var err error

		go func() {
			val, err = cvr.Recv("out")

			close(resCh)
		}()

		select {
		case <-resCh:
			if err != nil {
				log.Println(err)

				continue
			}

			log.Println("recv:", val)

		case <-time.After(readTimeout):
			log.Println("timeout")
		}
	}

	cancelRun()

	if err := <-runErrCh; err != nil {
		log.Println("conveyer finished with error:", err)

	} else {
		log.Println("conveyer finished")
	}
}
