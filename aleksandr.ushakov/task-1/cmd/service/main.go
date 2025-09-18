package main

import "fmt"

var (
	firstOperand  int
	secondOperand int
	operation     string
)

func main() {
	//Reading the first operand
	n, err := fmt.Scanln(&firstOperand)
	if err != nil || n != 1 {
		fmt.Println("Invalid first operand")
		return
	}

	//Reading the second operand
	n, err = fmt.Scanln(&secondOperand)
	if err != nil || n != 1 {
		fmt.Println("Invalid second operand")
		return
	}

	//Reading the operation
	n, err = fmt.Scanln(&operation)
	if err != nil || n != 1 {
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
