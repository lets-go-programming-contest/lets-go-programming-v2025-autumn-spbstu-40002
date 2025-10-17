package ktherr

import "errors"

var ErrPositionOutOfRange = errors.New("position out of range")

var ErrEmptyResult = errors.New("empty result")

var ErrInvalidItemCount = errors.New("item count must be between 1 and 10000")
