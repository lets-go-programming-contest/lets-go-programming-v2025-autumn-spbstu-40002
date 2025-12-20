package main

import (
	"fmt"
)

func main() {
	var firstOp, secondOp int
	var operation string

	if _, err := fmt.Scanln(&firstOp); err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	if _, err := fmt.Scanln(&secondOp); err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	if _, err := fmt.Scanln(&operation); err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operation {
	case "+":
		fmt.Println(firstOp + secondOp)
	case "-":
		fmt.Println(firstOp - secondOp)
	case "*":
		fmt.Println(firstOp * secondOp)
	case "/":
		if secondOp == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(firstOp / secondOp)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
