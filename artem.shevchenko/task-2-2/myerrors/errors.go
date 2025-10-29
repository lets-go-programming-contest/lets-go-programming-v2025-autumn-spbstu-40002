package myerrors

import (
	"errors"
)

var (
	ErrIncorrectDishesCount    = errors.New("the number of dishes must be between 1 and 10000")
	ErrIncorrectRaiting        = errors.New("the raiting of dish must be between -10000 and 10000")
	ErrIncorrectPickedDish     = errors.New("the picked dish must be between 1 and number of dishes")
	ErrIncorrectHeapPushedType = errors.New("the type being added to the heap is not an integer")
	ErrNothingToDelete         = errors.New("there are no elements in the heap to remove")
)
