package main

import "fmt"

func main() {
	var firstNumber int
	if _, err := fmt.Scan(&firstNumber); err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	var secondNumber int
	if _, err := fmt.Scan(&secondNumber); err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	var operator string
	if _, err := fmt.Scan(&operator); err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operator {
	case "+":
		fmt.Println(firstNumber + secondNumber)
	case "-":
		fmt.Println(firstNumber - secondNumber)
	case "*":
		fmt.Println(firstNumber * secondNumber)
	case "/":
		if secondNumber == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(firstNumber / secondNumber)
	default:
		fmt.Println("Invalid operation")
	}
}
