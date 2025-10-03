package main

import "fmt"

func main() {
	var (
		argOne int
		argTwo int
		symbol string
	)

	_, err := fmt.Scan(&argOne)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	_, err = fmt.Scan(&argTwo)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}
	_, err = fmt.Scan(&symbol)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch symbol {
	case "+":
		fmt.Println(argOne + argTwo)
	case "-":
		fmt.Println(argOne - argTwo)
	case "*":
		fmt.Println(argOne * argTwo)
	case "/":
		if argTwo == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(argOne / argTwo)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
