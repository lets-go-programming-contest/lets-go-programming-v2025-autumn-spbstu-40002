package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidOperator    = errors.New("invalid operator")
	ErrInvalidTemperature = errors.New("invalid temperature")
)

func parseBorder(border string) (string, int, error) {
	fields := strings.Fields(border)
	if len(fields) != 2 {
		return "", 0, ErrInvalidTemperature
	}
	if fields[0] != "<=" && fields[0] != ">=" {
		return "", 0, ErrInvalidOperator
	}
	temp, err := strconv.Atoi(fields[1])
	if err != nil {
		return "", 0, ErrInvalidTemperature
	}
	if temp < 15 || temp > 30 {
		return "", 0, ErrInvalidTemperature
	}
	return fields[0], temp, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func processEmployees(employees int) {
	minTemp := 15
	maxTemp := 30
	impossibleCondition := false
	for i := 0; i < employees; i++ {
		var oper string
		var temp string
		_, err := fmt.Scan(&oper, &temp)
		if err != nil {
			fmt.Println(-1)
			impossibleCondition = true
			continue
		}
		operation, temperature, err := parseBorder(fmt.Sprintf("%s %s", oper, temp))
		if err != nil {
			fmt.Println(-1)
			impossibleCondition = true
			continue
		}
		switch operation {
		case ">=":
			minTemp = max(minTemp, temperature)
		case "<=":
			maxTemp = min(maxTemp, temperature)
		}
		if minTemp > maxTemp {
			fmt.Println(-1)
			impossibleCondition = true
			continue
		}
		fmt.Println(minTemp)
	}
}

func main() {
	var departments int
	_, err := fmt.Scan(&departments)
	if err != nil {
		fmt.Println("Error reading number of departments:", err)
		return
	}
	for i := 0; i < departments; i++ {
		var employees int
		_, err := fmt.Scan(&employees)
		if err != nil {
			fmt.Println("Error reading number of employees:", err)
			return
		}
		processEmployees(employees)
	}
}
