package main

import "fmt"

func main() {
	var number1 int
	var number2 int
	var operator string

	_, err := fmt.Scan(&number1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&number2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&operator)
	if err != nil {
		return
	}

	switch operator {
	case "+":
		fmt.Println(number1 + number2)
	case "-":
		fmt.Println(number1 - number2)
	case "/":
		if number2 != 0 {
			fmt.Println(number1 / number2)
		} else {
			fmt.Println("Division by zero")
		}
	case "*":
		fmt.Println(number1 * number2)
	default:
		fmt.Println("Invalid operation")
	}
}
