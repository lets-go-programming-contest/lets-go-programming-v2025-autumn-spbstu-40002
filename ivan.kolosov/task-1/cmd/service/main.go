package main

import "fmt"

func main() {
	var firstNumber int

	_, err := fmt.Scan(&firstNumber)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	var secondNumber int
	_, err = fmt.Scan(&secondNumber)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	var operator string
	_, err = fmt.Scan(&operator)
	if err != nil {
		return
	}

	var answer int
	switch operator {
	case "+":
		answer = firstNumber + secondNumber
	case "-":
		answer = firstNumber - secondNumber
	case "*":
		answer = firstNumber * secondNumber
	case "/":
		if secondNumber == 0 {
			fmt.Println("Division by zero")
			return
		}
		answer = firstNumber / secondNumber
	default:
		fmt.Println("Invalid operation")
		return
	}
	fmt.Println(answer)
}
