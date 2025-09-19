package main

import (
	"fmt"
)

func main() {
	var firstVar, secondVar int
	var operation string

	_, err := fmt.Scanln(&firstVar)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scanln(&secondVar)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scanln(&operation)
	if err != nil {
		return
	}

	switch operation {
	case "+":
		fmt.Println(firstVar + secondVar)
	case "/":
		if secondVar == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(firstVar / secondVar)
	case "-":
		fmt.Println(firstVar - secondVar)
	case "*":
		fmt.Println(firstVar * secondVar)
	default:
		fmt.Println("Invalid operation")
		return
	}
}
