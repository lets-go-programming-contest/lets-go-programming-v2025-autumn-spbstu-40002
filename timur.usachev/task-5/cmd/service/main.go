package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/t1wt/task-5/pkg/conveyer"
	"github.com/t1wt/task-5/pkg/handlers"
)

func main() {
	cvr := conveyer.New(4)

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

	repCount := 3
	for i := 0; i < repCount; i++ {
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
			fmt.Println("recv:", val)
		case <-time.After(2 * time.Second):
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
