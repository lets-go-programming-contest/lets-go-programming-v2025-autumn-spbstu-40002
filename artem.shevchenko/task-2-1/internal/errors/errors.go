package errors

import "errors"

var (
	ErrIncorrectTemperature      = errors.New("the temperature is not between 15 and 30")
	ErrIncorrectOperator         = errors.New("the range sign can only be <= or >=")
	ErrIncorrectDepartmentsCount = errors.New("the number of departments must be between 1 and 1000")
	ErrIncorrectEmployeeCount    = errors.New("the number of employees must be between 1 and 1000")
)
