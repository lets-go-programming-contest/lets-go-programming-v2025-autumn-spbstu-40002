package inputError

import "errors"

var (
	ErrInvalidNumberDishes    = errors.New("The number of dishes on the table should be set from 1 to 10,000")
	ErrInvalidPriorityDish    = errors.New("The priority of the dish should be set from -10000 to 10000")
	ErrInvalidPreferredDishID = errors.New("The ID of the preferred dish should be no more than the number of dishes themselves.")
	ErrInvalidTypeData        = errors.New("Invalid data type")
)
