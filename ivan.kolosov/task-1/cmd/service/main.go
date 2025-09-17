package main

import "fmt"

func main() {
	var first_number int

	_, err := fmt.Scan(&first_number)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	var second_number int
	_, err = fmt.Scan(&second_number)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	var operator string
	fmt.Scan(&operator)

	var answer int
	switch operator {
	case "+":
		answer = first_number + second_number
	case "-":
		answer = first_number - second_number
	case "*":
		answer = first_number * second_number
	case "/":
		if second_number == 0 {
			fmt.Println("Division by zero")
			return
		}
		answer = first_number / second_number
	default:
		fmt.Println("Invalid operation")
		return
	}
	fmt.Println(answer)
}
