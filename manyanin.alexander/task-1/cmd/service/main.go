package main

import "fmt"

func main() {
	var first_operand int
	var second_operand int

	var operator string

	_, err := fmt.Scan(&first_operand)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&second_operand)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&operator)
	if err != nil {
		return
	}

	switch operator {
	case "*":
		fmt.Println(first_operand * second_operand)
	case "/":
		if second_operand == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(first_operand / second_operand)
	case "-":
		fmt.Println(first_operand - second_operand)
	case "+":
		fmt.Println(first_operand + second_operand)

	default:
		fmt.Println("Invalid operation")
	}
}
