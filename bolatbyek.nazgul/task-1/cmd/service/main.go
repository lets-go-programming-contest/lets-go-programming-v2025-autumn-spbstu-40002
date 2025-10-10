package main

import "fmt"

func main() {
	var firstOperand, secondOperand int
	var operation string

	//Reading the first operand
	_, err := fmt.Scanln(&firstOperand)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	//Reading the second operand
	_, err = fmt.Scanln(&secondOperand)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	//Reading the operation
	_, err = fmt.Scanln(&operation)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	//Calculating the expression
	var result int
	switch operation {
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

	//Writing the result
	fmt.Println(result)
}
