package main

import "fmt"

func main() {
	var a, b int
	var op string

	_, err := fmt.Scanln(&a)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scanln(&b)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scanln(&op)
	if err != nil {
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
		} else {
			fmt.Println(a / b)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
