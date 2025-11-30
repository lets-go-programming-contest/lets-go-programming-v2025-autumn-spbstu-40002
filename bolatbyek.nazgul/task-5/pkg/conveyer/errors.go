package conveyer

import "errors"

const chanNotFoundMsg = "chan not found"

var ErrChannelNotFound = errors.New(chanNotFoundMsg)
