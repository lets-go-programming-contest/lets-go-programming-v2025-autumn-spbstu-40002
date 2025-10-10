package service

import (
	"fmt"
)

func main() {
	var firstOp, secondOp int
	var operation string

	_, err := fmt.Scan(&firstOp)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&secondOp)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&operation)
	if err != nil {
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
