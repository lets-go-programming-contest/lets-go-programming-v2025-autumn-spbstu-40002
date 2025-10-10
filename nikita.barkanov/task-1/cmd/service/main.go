package main

import "fmt"

func main() {
	var (
		num1, num2 int
		operator   string
		ans        int
	)

	_, err := fmt.Scan(&num1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&num2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&operator)
	if err != nil {
		return
	}

	switch operator {
	case "+":
		ans = num1 + num2
	case "-":
		ans = num1 - num2
	case "*":
		ans = num1 * num2
	case "/":
		if num2 == 0 {
			fmt.Println("Division by zero")
			return
		}
		ans = num1 / num2
	default:
		fmt.Println("Invalid operation")
		return
	}

	fmt.Println(ans)
}
