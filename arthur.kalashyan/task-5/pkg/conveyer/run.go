package conveyer

import (
	"context"
	"errors"
	"sync"
)

var errChanNotFound = errors.New("chan not found")

func (c *conveyor) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	errCh := make(chan error, 1)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}
