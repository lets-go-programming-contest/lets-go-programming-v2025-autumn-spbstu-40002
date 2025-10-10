package main

import "fmt"

func calculate(a int, b int, operand rune) interface{} {
	switch operand {
	case '/':
		if b == 0 {
			return "Division by zero"
		}
		return a / b
	case '+':
		return a + b
	case '-':
		return a - b
	case '*':
		return a * b
	default:
		return "Invalid operation"
	}
}

func main() {
	var a int
	var err error
	_, err = fmt.Scanln(&a)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	var b int
	_, err = fmt.Scanln(&b)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	var op rune
	_, err = fmt.Scanf("%c\n", &op)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	fmt.Println(calculate(a, b, op))
}
