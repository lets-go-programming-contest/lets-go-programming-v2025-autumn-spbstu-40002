package main

import "fmt"

const (
	minDepartNumber = 1
	maxDepartNumber = 1000
	minWorkersNum   = 1
	maxWorkersNum   = 1000
	minTemperature  = 15
	maxTemperature  = 30
)

type departTemperatureHandler struct {
	optimalTemperature int

	upperBound int
	lowerBound int // минимально допустимая температура, то есть начиная с нее температура приемлем

}

func (object *departTemperatureHandler) setTemperature(operator string, value int) error {
	switch operator {
	case ">=":
		object.lowerBound = value
	case "<=":
		object.upperBound = value
	default:
		return fmt.Errorf("Unknown operator: %s", operator)
	}
	return nil
}

func NewDepartTemperatureHandler(lBound int, uBound int) *departTemperatureHandler {
	return &departTemperatureHandler{
		optimalTemperature: lBound,

		lowerBound: lBound,
		upperBound: uBound,
	}
}

func main() {

	var departNumber int

	var _, err = fmt.Scanln(&departNumber)
	if (err != nil) || (departNumber > maxDepartNumber || departNumber < minDepartNumber) {
		fmt.Println("Invalid department number")
	}

	for i := 0; i < departNumber; i++ {
		var workersNum int
		_, err = fmt.Scanln(&workersNum)
		if (err != nil) || (workersNum > maxWorkersNum || workersNum < minWorkersNum) {
			fmt.Println("Invalid temperature value")
		}

		for j := 0; j < workersNum; j++ {
			var operator string
			var value int

			var handler = NewDepartTemperatureHandler(minTemperature, maxTemperature)

			_, err = fmt.Scanln(&operator)
			if err != nil {
				fmt.Println("Invalid operator")
			}
			_, err = fmt.Scanln(&value)
			if err != nil {
				fmt.Println("Incorrect value")
			}

			handler.setTemperature(operator, value)
		}

	}

}
