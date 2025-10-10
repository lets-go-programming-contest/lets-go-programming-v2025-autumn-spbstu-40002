package main

import (
	"fmt"
)

func main() {
	fmt.Println("Enter the first operand:")
	var a int
	_, err := fmt.Scan(&a)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	fmt.Println("Enter the second operand:")
	var b int
	_, err = fmt.Scan(&b)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	fmt.Println("Enter the operation (+, -, *, /):")
	var operation string
	_, err = fmt.Scan(&operation)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operation {
	case "+":
		fmt.Println("Result:", a+b)
	case "-":
		fmt.Println("Result:", a-b)
	case "*":
		fmt.Println("Result:", a*b)
	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println("Result:", a/b)
	default:
		fmt.Println("Invalid operation")
	}
}
