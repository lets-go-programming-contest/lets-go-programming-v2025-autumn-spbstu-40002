package main

import (
	"fmt"
)

func parse(a any, errText string) bool {
	_, err := fmt.Scanln(a)
	if err != nil {
		fmt.Println(errText)
		var errBuf string
		fmt.Scanln(&errBuf)
	}
	return err == nil
}

func calc(first, second int, operation string) {
	switch operation {
	case "+":
		fmt.Println(first + second)
	case "-":
		fmt.Println(first - second)
	case "/":
		if second != 0 {
			fmt.Println(first / second)
		} else {
			fmt.Println("Division by zero")
		}
	case "*":
		fmt.Println(first * second)
	default:
		fmt.Println("Invalid operation")
	}
}

func main() {
	var (
		firstOperand  int
		secondOperand int
		operation     string
	)

	if !parse(&firstOperand, "Invalid first operand") ||
		!parse(&secondOperand, "Invalid second operand") ||
		!parse(&operation, "Invalid operation") {
		return
	}

	calc(firstOperand, secondOperand, operation)
}
