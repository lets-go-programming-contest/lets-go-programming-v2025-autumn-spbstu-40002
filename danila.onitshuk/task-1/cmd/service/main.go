package main

import "fmt"

var (
	num1     int
	num2     int
	operator string
)

func main() {
	_, err1 := fmt.Scan(&num1)
	_, err2 := fmt.Scan(&num2)
	_, err3 := fmt.Scan(&operator)

	switch {
	case err1 != nil:
		fmt.Println("Invalid first operand")
	case err2 != nil:
		fmt.Println("Invalid second operand")
	case err3 != nil:
		return
	case operator == "+":
		fmt.Println(num1 + num2)
	case operator == "-":
		fmt.Println(num1 - num2)
	case operator == "*":
		fmt.Println(num1 * num2)
	case operator == "/":
		if num2 == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(num1 / num2)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
