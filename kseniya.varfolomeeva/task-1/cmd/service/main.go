package main

import (
	"fmt"
)

func main() {
	fmt.Println("Enter the first operand:")
	var a int
	_, err := fmt.Scan(&a)
	if err != nil {
		fmt.Println("Error: invalid input for first operand")
		return
	}

	fmt.Println("Enter the second operand:")
	var b int
	_, err = fmt.Scan(&b)
	if err != nil {
		fmt.Println("Error: invalid input for second operand")
		return
	}

	fmt.Println("Enter the operation (+, -, *, /):")
	var operation string
	_, err = fmt.Scan(&operation)
	if err != nil {
		fmt.Println("Error: invalid operation")
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
			fmt.Println("Error: division by zero")
			return
		}
		fmt.Println("Result:", a/b)
	default:
		fmt.Println("Error: unknown operation")
	}
}
