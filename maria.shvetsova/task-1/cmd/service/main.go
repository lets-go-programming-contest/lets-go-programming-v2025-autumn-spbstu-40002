package main

import (
	"fmt"
)

func main() {
	var a, b int
	var op string

	_, errScanFirstOperand := fmt.Scan(&a)
	if errScanFirstOperand != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, errScanSecondOperand := fmt.Scan(&b)
	if errScanSecondOperand != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, errScanOperator := fmt.Scan(&op)
	if errScanOperator != nil {
		return
	}
	isValid := (op == "+" || op == "-" || op == "*" || op == "/")
	if !isValid {
		fmt.Println("Invalid operation")
		return
	}

	switch op {
	case "+":
		fmt.Println(a + b)
	case "-":
		fmt.Println(a - b)
	case "*":
		fmt.Println(a * b)
	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(a / b)
	}
}
