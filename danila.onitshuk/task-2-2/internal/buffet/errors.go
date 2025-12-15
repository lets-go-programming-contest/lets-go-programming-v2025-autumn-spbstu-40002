package buffet

import "errors"

var (
	ErrInvalidNumberDishes = errors.New(
		"the number of dishes on the table should be set from 1 to 10,000",
	)
	ErrInvalidPriorityDish = errors.New(
		"the priority of the dish should be set from -10000 to 10000",
	)
	ErrInvalidPreferredDishID = errors.New(
		"the ID of the preferred dish should be no more than the number of dishes themselves",
	)
	ErrInvalidTypeData = errors.New("invalid data type")
	ErrUnderFlow       = errors.New("under flow")
)
