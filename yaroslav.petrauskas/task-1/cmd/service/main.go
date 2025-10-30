package main

import "fmt"

func main() {
	var a, b int
	var op string

	fmt.Scan(&a, &b, &op)

	switch {
	case op == "+":
		fmt.Println(a + b)
	case op == "-":
		fmt.Println(a - b)
	case op == "*":
		fmt.Println(a * b)
	case op == "/" && b == 0:
		fmt.Println("Division by zero")
	case op == "/":
		fmt.Println(a / b)
	default:
		fmt.Println("Invalid operation")
	}
}
