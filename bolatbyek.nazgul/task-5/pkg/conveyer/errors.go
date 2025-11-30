package conveyer

import "errors"

var chanNotFoundMsg = "chan not found"

var ErrChannelNotFound = errors.New(chanNotFoundMsg)
