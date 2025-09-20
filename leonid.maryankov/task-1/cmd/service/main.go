package main

import "fmt"

func main() {
	var (
		argOne float64
		argTwo float64
		symbol string
	)

	_, errOne := fmt.Scan(&argOne)
	if errOne != nil {
		fmt.Println("Invalid first operand")
		return
	}
	_, errTwo := fmt.Scan(&argTwo)
	if errTwo != nil {
		fmt.Println("Invalid second operand")
		return
	}
	_, errSym := fmt.Scan(&symbol)
	if errSym != nil {
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
