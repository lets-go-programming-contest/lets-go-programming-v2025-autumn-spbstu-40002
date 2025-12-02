package ktherr

import "errors"

var ErrPositionOutOfRange = errors.New("position out of range")

var ErrEmptyResult = errors.New("empty result")

var ErrInvalidItemCount = errors.New("item count must be between 1 and 10000")

var ErrReadItemCount = errors.New("failed to read item count")

var ErrReadValue = errors.New("failed to read value")

var ErrValueOutOfRange = errors.New("value must be between -10000 and 10000")

var ErrReadPosition = errors.New("failed to read position")
