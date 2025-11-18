package handlers

import "errors"

const (
	noDecoratorData        = "no decorator"
	textForDecoratorString = "decorated: "
	noMultiplexerData      = "no multiplexer"
)

var (
	ErrNoDecorator = errors.New("can’t be decorated")
)
