package main

import "fmt"

func main() {
	var firstOperand int
	var secondOperand int
	var operator string
	var answer int

	_, err := fmt.Scan(&firstOperand)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&secondOperand)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&operator)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operator {
	case "+":
		answer = firstOperand + secondOperand
	case "-":
		answer = firstOperand - secondOperand
	case "*":
		answer = firstOperand * secondOperand
	case "/":
		if secondOperand == 0 {
			fmt.Println("Division by zero")
			return
		}
		answer = firstOperand / secondOperand
	default:
		fmt.Println("Invalid operation")
		return
	}
	fmt.Println(answer)
}
