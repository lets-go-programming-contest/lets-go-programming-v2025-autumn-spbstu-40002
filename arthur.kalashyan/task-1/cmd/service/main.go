package main

import (
	"fmt"
)

func main() {
	var firstOperand, secondOperand int
	var operator string

	if _, err := fmt.Scanln(&firstOperand); err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	if _, err := fmt.Scanln(&secondOperand); err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	if _, err := fmt.Scanln(&operator); err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operator {
	case "+":
		fmt.Println(firstOperand + secondOperand)
	case "-":
		fmt.Println(firstOperand - secondOperand)
	case "*":
		fmt.Println(firstOperand * secondOperand)
	case "/":
		if secondOperand == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(firstOperand / secondOperand)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
