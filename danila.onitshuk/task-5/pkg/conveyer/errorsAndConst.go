package conveyer

import "errors"

const (
	undefinedData = "undefined"
)

var (
	ErrNoChannel = errors.New("chan not found")
)
