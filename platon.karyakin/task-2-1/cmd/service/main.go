package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	defaultMinTemp = 15
	defaultMaxTemp = 30
)

type TemperatureController struct {
	min int
	max int
}

func NewTemperatureController() *TemperatureController {
	return &TemperatureController{
		min: defaultMinTemp,
		max: defaultMaxTemp,
	}
}

func (tc *TemperatureController) ApplyConstraint(operation string, temperature int) {
	switch operation {
	case ">=":
		if temperature > tc.min {
			tc.min = temperature
		}
	case "<=":
		if temperature < tc.max {
			tc.max = temperature
		}
	}
}

func (tc *TemperatureController) GetCurrentTemp() string {
	if tc.min <= tc.max {
		return strconv.Itoa(tc.min)
	}

	return "-1"
}

type Department struct {
	workerCount int
	controller  *TemperatureController
}

func NewDepartment(workerCount int) *Department {
	return &Department{
		workerCount: workerCount,
		controller:  NewTemperatureController(),
	}
}

func (d *Department) ProcessWorkerRequirement(operation string, temperature int) string {
	d.controller.ApplyConstraint(operation, temperature)

	return d.controller.GetCurrentTemp()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	depCount, _ := strconv.Atoi(scanner.Text())

	for range depCount {
		scanner.Scan()

		workersNumber, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return
		}

		department := NewDepartment(workersNumber)

		for range workersNumber {
			scanner.Scan()
			line := scanner.Text()
			parts := strings.Fields(line)
			operation := parts[0]

			temp, err := strconv.Atoi(parts[1])
			if err != nil {
				return
			}

			result := department.ProcessWorkerRequirement(operation, temp)
			fmt.Println(result)
		}
	}
}
