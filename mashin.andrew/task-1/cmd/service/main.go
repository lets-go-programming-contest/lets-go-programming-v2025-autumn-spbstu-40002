package main

import "fmt"

func main() {
	var number1 int
	_, err1 := fmt.Scan(&number1)

	if err1 != nil {
		fmt.Printf("Invalid first operand\n")
		return
	}

	var number2 int
	_, err2 := fmt.Scan(&number2)

	if err2 != nil {
		fmt.Printf("Invalid second operand\n")
		return
	}

	var operator string
	if _, err := fmt.Scan(&operator); err != nil {
		return
	}

	switch operator {
	case "+":
		fmt.Print(number1 + number2, "\n")
	case "-":
		fmt.Print(number1 - number2, "\n")
	case "/":
		if number2 != 0 {
			fmt.Print(number1 / number2, "\n")
		} else {
			fmt.Printf("Division by zero\n")
		}
	case "*":
		fmt.Print(number1 * number2, "\n")
	default:
		fmt.Printf("Invalid operation\n")
	}
}
