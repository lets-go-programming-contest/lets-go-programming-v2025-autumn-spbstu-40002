package handlers

import "errors"

var (
	ErrNoDecorator = errors.New("can't be decorated")
)

const (
	noDecoratorData        = "no decorator"
	textForDecoratorString = "decorated: "
	noMultiplexerData      = "no multiplexer"
)
