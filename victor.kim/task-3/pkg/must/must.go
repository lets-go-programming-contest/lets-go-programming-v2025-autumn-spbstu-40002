package must

import (
	"fmt"
	"io"
)

func Must(op string, err error) {
	if err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}
}

func Close[T io.Closer](path string, c T) {
	Must(fmt.Sprintf("close %q", path), c.Close())
}
