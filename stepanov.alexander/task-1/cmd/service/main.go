package main

import (
	"fmt"
)

func main() {
	var a, b float32
	var operation string

	_, err := fmt.Scan(&a)
	if err != nil {
		fmt.Println("Invalid first operand")
	}
	_, er := fmt.Scan(&b)
	if er != nil {
		fmt.Println("Invalid second operand")
	}

	fmt.Scan(&operation)

	switch operation {
	case "+":
		fmt.Println(a + b)
	case "-":
		fmt.Println(a - b)
	case "*":
		fmt.Println(a * b)
	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(a / b)
		}
	default:
		fmt.Println("Invalid operation")
	}

}
