package main

import (
	"errors"
	"fmt"
)

const (
	minDepartNumber = 1
	maxDepartNumber = 1000
	minWorkersNum   = 1
	maxWorkersNum   = 1000
	minTemperature  = 15
	maxTemperature  = 30
)

var ErrUnknownOperator = errors.New("unknown operator")

type DepartTemperatureHandler struct {
	optimalTemperature int

	upperBound int
	lowerBound int // минимально допустимая температура, то есть начиная с нее температура приемлем
}

func (object *DepartTemperatureHandler) setTemperature(operator string, value int) error {
	switch operator {
	case ">=":
		object.lowerBound = value
		if object.optimalTemperature < object.lowerBound {
			object.optimalTemperature = object.lowerBound
		}
	case "<=":
		object.upperBound = value
		if object.optimalTemperature > object.upperBound {
			object.optimalTemperature = object.upperBound
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnknownOperator, operator)
	}

	if object.upperBound < object.lowerBound {
		object.optimalTemperature = -1
	}

	return nil
}

func (object *DepartTemperatureHandler) getTemperature() int {
	return object.optimalTemperature
}

func NewDepartTemperatureHandler(lBound int, uBound int) *DepartTemperatureHandler {
	return &DepartTemperatureHandler{
		optimalTemperature: lBound,

		lowerBound: lBound,
		upperBound: uBound,
	}
}

func main() {
	var departNumber int

	_, err := fmt.Scanln(&departNumber)

	if err != nil || departNumber > maxDepartNumber || departNumber < minDepartNumber {
		fmt.Println("Invalid department number")

		return
	}

	for range make([]struct{}, departNumber) {
		var workersNum int

		_, err = fmt.Scanln(&workersNum)

		if err != nil || workersNum > maxWorkersNum || workersNum < minWorkersNum {
			fmt.Println("Invalid temperature value")

			return
		}

		handler := NewDepartTemperatureHandler(minTemperature, maxTemperature)
		for _ = range make([]struct{}, workersNum) {
			var operator string

			var value int

			_, err = fmt.Scanln(&operator, &value)
			if err != nil {
				fmt.Println("Invalid input:", err)

				continue
			}

			if err := handler.setTemperature(operator, value); err != nil {
				fmt.Println("Error:", err)

				continue
			}

			temp := handler.getTemperature()

			fmt.Println(temp)
		}
	}
}
