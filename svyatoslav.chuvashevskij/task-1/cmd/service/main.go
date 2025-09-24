package main

import "fmt"

func main() {
	var firstOperand int
	var secondOperand int
	var operationSymbol string
	var result int

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

	_, err = fmt.Scan(&operationSymbol)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operationSymbol {
	case "+":
		result = firstOperand + secondOperand
	case "-":
		result = firstOperand - secondOperand
	case "*":
		result = firstOperand * secondOperand
	case "/":
		if secondOperand == 0 {
			fmt.Println("Division by zero")
			return
		}
		result = firstOperand / secondOperand
	default:
		fmt.Println("Invalid operation")
		return
	}

	fmt.Println(result)
}
