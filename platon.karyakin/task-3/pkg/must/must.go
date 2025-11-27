package must

import (
	"fmt"
	"io"
)

func Must(operation string, err error) {
	if err != nil {
		panic(fmt.Errorf("operation %q failed: %w", operation, err))
	}
}

func Close[T io.Closer](path string, c T) {
	Must(fmt.Sprintf("close %q", path), c.Close())
}
