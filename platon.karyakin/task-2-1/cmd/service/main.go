package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

const (
	initialMinTemp = 15
	initialMaxTemp = 30
)

var errFormat = errors.New("invalid input format or value")

func readOperationTemp() (string, int, error) {
	var (
		operation string
		temp      int
	)

	_, err := fmt.Fscan(os.Stdin, &operation, &temp)
	if err != nil {
		return "", 0, errFormat
	}

	if operation != ">=" && operation != "<=" {
		return "", 0, errFormat
	}

	if temp < initialMinTemp || temp > initialMaxTemp {
		return "", 0, errFormat
	}

	return operation, temp, nil
}

type TemperatureManager struct {
	minTemp int
	maxTemp int
}

func NewTemperatureManager(minInit, maxInit int) *TemperatureManager {
	return &TemperatureManager{
		minTemp: minInit,
		maxTemp: maxInit,
	}
}

func (tm *TemperatureManager) Update(operation string, temperature int) error {
	if operation != ">=" && operation != "<=" {
		return errFormat
	}

	if temperature < initialMinTemp || temperature > initialMaxTemp {
		return errFormat
	}

	if operation == ">=" {
		if temperature > tm.minTemp {
			tm.minTemp = temperature
		}
	} else if operation == "<=" {
		if temperature < tm.maxTemp {
			tm.maxTemp = temperature
		}
	}

	return nil
}

func (tm *TemperatureManager) Get() (int, int) {
	return tm.minTemp, tm.maxTemp
}

func main() {
	var departmentsCount int

	if _, err := fmt.Fscan(os.Stdin, &departmentsCount); err != nil {
		return
	}

	for range departmentsCount {
		var employeesCount int

		if _, err := fmt.Fscan(os.Stdin, &employeesCount); err != nil {
			return
		}

		currentMinTemp := initialMinTemp
		currentMaxTemp := initialMaxTemp

		tempManager := NewTemperatureManager(currentMinTemp, currentMaxTemp)

		for range employeesCount {
			operation, temperature, err := readOperationTemp()
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			if err := tempManager.Update(operation, temperature); err != nil {
				log.Fatalf("error: %v", err)
			}

			currentMinTemp, currentMaxTemp = tempManager.Get()

			if currentMinTemp <= currentMaxTemp {
				fmt.Println(currentMinTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
