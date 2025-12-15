package main

import "fmt"

func main() {
	var a, b int
	var op string

	if _, e := fmt.Scanln(&a); e != nil {
		fmt.Println("Invalid first operand")
		return
	}

	if _, e := fmt.Scanln(&b); e != nil {
		fmt.Println("Invalid second operand")
		return
	}

	if _, e := fmt.Scanln(&op); e != nil {
		fmt.Println("Invalid operation")
		return
	}

	var res int
	switch op {
	case "+":
		res = a + b
	case "-":
		res = a - b
	case "*":
		res = a * b
	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
			return
		}
		res = a / b
	default:
		fmt.Println("Invalid operation")
		return
	}

	fmt.Println(res)
}
