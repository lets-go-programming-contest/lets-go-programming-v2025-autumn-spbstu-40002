package main

import "fmt"

func main() {
	var a, b int
	var op string

	if _, err := fmt.Scan(&a); err != nil {
		return
	}

	if _, err := fmt.Scan(&b); err != nil {
		return
	}

	if _, err := fmt.Scan(&op); err != nil {
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
			return
		}
		fmt.Println(a / b)
	default:
		return
	}
}
