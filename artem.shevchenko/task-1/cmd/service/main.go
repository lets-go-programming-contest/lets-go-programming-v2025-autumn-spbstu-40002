package main

import "fmt"

func main() {
	var num1, num2 int
	var operator string

	// Read the first operand
	_, err := fmt.Scan(&num1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	// Read the second operand
	_, err = fmt.Scan(&num2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	// Read the operation sign
	_, err = fmt.Scan(&operator)
	if err != nil {
		return
	}

	// Сalculate the expression depending on the operator
	switch operator {
	case "+":
		fmt.Println(num1 + num2)
	case "-":
		fmt.Println(num1 - num2)
	case "*":
		fmt.Println(num1 * num2)
	case "/":
		if num2 == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(num1 / num2)
	default:
		fmt.Println("Invalid operation")
		return
	}
}
