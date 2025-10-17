package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var a, b int
	var op string

	// Read first operand
	if _, err := fmt.Fscan(in, &a); err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	// Read second operand
	if _, err := fmt.Fscan(in, &b); err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	// Read operator
	if _, err := fmt.Fscan(in, &op); err != nil {
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
	default:
		fmt.Println("Invalid operation")
	}
}
