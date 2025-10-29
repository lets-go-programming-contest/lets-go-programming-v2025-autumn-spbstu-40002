package main

import "fmt"

func main() {
	var firstОperand int
	var secondOperand int

	var operator string

	_, err := fmt.Scan(&firstОperand)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&secondOperand)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&operator)
	if err != nil {
		return
	}

	switch operator {
	case "*":
		fmt.Println(firstОperand * secondOperand)
	case "/":
		if secondOperand == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(firstОperand / secondOperand)
	case "-":
		fmt.Println(firstОperand - secondOperand)
	case "+":
		fmt.Println(firstОperand + secondOperand)

	default:
		fmt.Println("Invalid operation")
	}
}
