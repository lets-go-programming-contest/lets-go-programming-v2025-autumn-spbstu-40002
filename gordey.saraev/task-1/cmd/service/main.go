package main

import (
	"fmt"
)

func main() {
	var firstVar, operation, secondVar string

	fmt.Scanln(&firstVar)
	fmt.Scanln(&operation)
	fmt.Scanln(&secondVar)

	first, err := toInt(firstVar)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	second, err := toInt(secondVar)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	switch operation {
	case "+":
		fmt.Println(first + second)
	case "/":
		if second == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(first / second)
	case "-":
		fmt.Println(first - second)
	case "*":
		fmt.Println(first * second)
	default:
		fmt.Println("Invalid operation")
		return
	}

}

func toInt(s string) (int, error) {
	var num int
	_, err := fmt.Sscanf(s, "%d", &num)
	return num, err
}
