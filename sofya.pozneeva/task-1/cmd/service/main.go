package main

import "fmt"

func calculate(a int, b int, operand rune) interface{} {
	if operand == '/' && b == 0 {
		return "Division by zero"
	}
	switch operand {
	case '/':
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
		fmt.Print("Invalid first operand")
		return
	}

	var b int
	_, err = fmt.Scanln(&b)
	if err != nil {
		fmt.Print("Invalid second operand")
		return
	}

	var op rune
	_, err = fmt.Scanf("%c\n", &op)
	if err != nil {
		fmt.Print("Invalid operation")
		return
	}

	fmt.Println(calculate(a, b, op))
}
