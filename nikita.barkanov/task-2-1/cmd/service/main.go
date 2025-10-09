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

	lowerBound int // минимально допустимая температура, то есть начиная с нее температура приемлем
	upperBound int // максимально допустимая температура
}

func (object *DepartTemperatureHandler) setTemperature(operator string, value int) error {
	switch operator {
	case ">=":
		if object.lowerBound < value {
			object.lowerBound = value
		}

		if object.lowerBound < minTemperature {
			object.lowerBound = minTemperature
		}

		if object.optimalTemperature < object.lowerBound {
			object.optimalTemperature = object.lowerBound
		}
	case "<=":
		if object.upperBound > value {
			object.upperBound = value
		}

		if object.upperBound > maxTemperature {
			object.upperBound = maxTemperature
		}

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
	// Другие варианты циклов не работают
	// for i:= 0; i < N; i++ не пропускается линтером
	// for i := range departNumber не работает, потому что у меня не используются i и j в
	//	соответствующих циклах, отсюда лезет ошибка unused variable
	// конструкции типа for range departNumber так же не работают
	index1 := 0
	for index1 < departNumber {
		var workersNum int

		_, err = fmt.Scanln(&workersNum)

		if err != nil || workersNum > maxWorkersNum || workersNum < minWorkersNum {
			fmt.Println("Invalid workers num")

			return
		}

		handler := NewDepartTemperatureHandler(minTemperature, maxTemperature)

		index2 := 0
		for index2 < workersNum {
			var operator string

			var value int

			_, err = fmt.Scanln(&operator, &value)
			if err != nil {
				fmt.Println("Invalid input:", err)

				return
			}

			if err := handler.setTemperature(operator, value); err != nil {
				fmt.Println("Error:", err)

				return
			}

			temp := handler.getTemperature()

			fmt.Println(temp)

			index2++
		}

		index1++
	}
}
