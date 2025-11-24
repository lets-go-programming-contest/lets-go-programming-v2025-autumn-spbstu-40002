package conveyer

import "context"

func (c *conveyor) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.ensureChan(input)
	c.ensureChan(output)
	// append to internal registry (implement as slice fields)
}

func (c *conveyor) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	for _, in := range inputs {
		c.ensureChan(in)
	}
	c.ensureChan(output)
	// append to internal registry
}

func (c *conveyor) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.ensureChan(input)
	for _, out := range outputs {
		c.ensureChan(out)
	}
	// append to internal registry
}
