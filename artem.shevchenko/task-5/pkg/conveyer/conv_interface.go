package conveyer

import "context"

type conveyer interface {
	RegisterDecorator(
		fn func(
			ctx context.Context,
			input chan string,
			output chan string,
		) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		fn func(
			ctx context.Context,
			inputs []chan string,
			output chan string,
		) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		fn func(
			ctx context.Context,
			input chan string,
			outputs []chan string,
		) error,
		input string,
		outputs []string,
	)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}
