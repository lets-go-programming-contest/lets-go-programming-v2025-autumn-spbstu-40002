package department

import "errors"

var (
	ErrInvalidTypeData = errors.New("invalid data type")
	ErrCountDeparments = errors.New("the number of departments should be from 1 to 1000")
	ErrcountEmployees  = errors.New("the number of employees should be from 1 to 1000")
	ErrInvalidSign     = errors.New("invalid sign")
)
