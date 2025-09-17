package main

import "fmt"

func main() {
	var number1 int
	_, err1 := fmt.Scan(&number1)

	if err1 != nil {
		fmt.Println("Invalid first operand")
		return
	}

	var number2 int
	_, err2 := fmt.Scan(&number2)

	if err2 != nil {
		fmt.Println("Invalid second operand")
		return
	}

	var operator string
	if _, err := fmt.Scan(&operator); err != nil {
		return
	}

	switch operator {
	case "+":
		fmt.Print(number1 + number2)
	case "-":
		fmt.Print(number1 - number2)
	case "/":
		if number2 != 0 {
			fmt.Print(number1 / number2)
		} else {
			fmt.Print("Division by zero")
		}
	case "*":
		fmt.Print(number1 * number2)
	default:
		fmt.Print("Invalid operation")
	}
}
