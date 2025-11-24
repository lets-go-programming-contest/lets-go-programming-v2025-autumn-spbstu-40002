package conveyer

import (
	"errors"
)

func (c *conveyor) Send(input string, data string) error {
	ch, ok := c.getChan(input)
	if !ok {
		return errors.New("chan not found")
	}
	ch <- data
	return nil
}

func (c *conveyor) Recv(output string) (string, error) {
	ch, ok := c.getChan(output)
	if !ok {
		return "", errors.New("chan not found")
	}
	v, ok := <-ch
	if !ok {
		return Undefined, nil
	}
	return v, nil
}
