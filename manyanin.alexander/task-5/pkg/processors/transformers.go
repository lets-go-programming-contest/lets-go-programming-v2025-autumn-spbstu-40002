package processors

import (
	"context"
	"fmt"
	"strings"
)

var cannotTransform = fmt.Errorf("transformation not possible")

const (
	modifiedPrefix = "processed: "
	skipTransform  = "skip_transform"
)

func AddPrefix(
	ctx context.Context,
	inCh chan string,
	outCh chan string,
) error {
	defer func() {
		close(outCh)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-inCh:
			if !ok {
				return nil
			}

			if strings.Contains(data, skipTransform) {
				return cannotTransform
			}

			var result string
			if strings.HasPrefix(data, modifiedPrefix) {
				result = data
			} else {
				result = modifiedPrefix + data
			}

			select {
			case <-ctx.Done():
				return nil
			case outCh <- result:
			}
		}
	}
}

func ToUpperCase(
	ctx context.Context,
	inCh chan string,
	outCh chan string,
) error {
	defer close(outCh)

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-inCh:
			if !ok {
				return nil
			}

			result := strings.ToUpper(data)

			select {
			case <-ctx.Done():
				return nil
			case outCh <- result:
			}
		}
	}
}
