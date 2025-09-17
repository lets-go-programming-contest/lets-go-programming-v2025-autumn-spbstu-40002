package main

import "fmt"

func main() {
	var num1 int
	var num2 int
	var operator string

	// Read the first operand
	fmt.Println("Enter the first number: ")
	_, err := fmt.Scan(&num1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	// Read the second operand
	fmt.Println("Enter the second number:")
	_, err = fmt.Scan(&num2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	// Read the operation sign
	fmt.Println("Enter operator: ")
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
