package main

import (
	"fmt"
)

func main() {
	var firstNum int
	var secondNum int
	var op string

	if _, err := fmt.Scanln(&firstNum); err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	if _, err := fmt.Scanln(&secondNum); err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	if _, err := fmt.Scanln(&op); err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch op {
	case "+":
		fmt.Println(firstNum + secondNum)
	case "-":
		fmt.Println(firstNum - secondNum)
	case "*":
		fmt.Println(firstNum * secondNum)
	case "/":
		if secondNum == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(firstNum / secondNum)
	default:
		fmt.Println("Invalid operation")
	}
}
