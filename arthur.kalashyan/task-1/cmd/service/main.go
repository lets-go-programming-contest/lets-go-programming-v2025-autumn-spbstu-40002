package main

import (
	"fmt"
)

func main() {
	var first, second int
	var op string

	if _, err := fmt.Scanln(&first); err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	if _, err := fmt.Scanln(&second); err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	if _, err := fmt.Scanln(&op); err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch op {
	case "+":
		fmt.Println(first + second)
	case "-":
		fmt.Println(first - second)
	case "*":
		fmt.Println(first * second)
	case "/":
		if second == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(first / second)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
