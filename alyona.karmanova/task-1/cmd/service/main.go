package main

import (
	"fmt"
)

func main() {
	var a, b float64
	_, err := fmt.Scan(&a)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&b)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	var action string
	_, err = fmt.Scan(&action)
	if err != nil {
		return
	}

	switch action {
	case "+":
		fmt.Println(a + b)
	case "-":
		fmt.Println(a - b)
	case "*":
		fmt.Println(a * b)
	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(a / b)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
