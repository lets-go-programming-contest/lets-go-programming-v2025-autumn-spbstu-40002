package conveyer

import (
	"context"
	"sync"
)

type Conveyer interface {
	RegisterDecorator(
		fn func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		fn func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		fn func(ctx context.Context, input chan string, outputs []chan string) error,
		input string,
		outputs []string,
	)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type conveyer struct {
	mu       sync.RWMutex
	channels map[string]chan string
	handlers []func(ctx context.Context) error
	size     int
}

func (c *conveyer) Send(input string, data string) error {
	return nil
}

func (c *conveyer) Recv(output string) (string, error) {
	return "", nil
}

func (c *conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	input, output string,
) {

}

func (c *conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string, output string,
) {
}

func (c *conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string, outputs []string,
) {
}

func (c *conveyer) Run(ctx context.Context) error {
	return nil
}

func New(size int) Conveyer {
	return &conveyer{
		channels: make(map[string]chan string),
		size:     size,
	}
}
