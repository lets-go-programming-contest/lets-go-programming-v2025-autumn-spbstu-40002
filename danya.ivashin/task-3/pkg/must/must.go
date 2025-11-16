package must

import (
	"fmt"
	"io"
)

func Must(op string, err error) {
	if err != nil {
		panic(fmt.Sprintf("op: %s", err))
	}
}

func Close[T io.Closer](path string, closer T) {
	Must(fmt.Sprintf("close %q", path), closer.Close())
}
