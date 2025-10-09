package main

import (
	"errors"
	"fmt"
)

var errDivisionByZero = errors.New("Division by zero")

func plus(x, y int) (int, error) {
	return x + y, nil
}

func minus(x, y int) (int, error) {
	return x - y, nil
}

func mult(x, y int) (int, error) {
	return x * y, nil
}

func div(x, y int) (int, error) {
	if y == 0 {
		return 0, errDivisionByZero
	}
	return x / y, nil
}

func main() {
	mapper := map[string]func(x, y int) (int, error){
		"+": plus,
		"-": minus,
		"*": mult,
		"/": div,
	}
	var (
		x, y int
		op   string
	)
	if _, err := fmt.Scan(&x); err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	if _, err := fmt.Scan(&y); err != nil {
		fmt.Println("Invalid second operand")
		return
	}
	if _, err := fmt.Scan(&op); err != nil {
		fmt.Println("Invalid operation")
		return
	}
	if f, ok := mapper[op]; ok {
		res, err := f(x, y)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(res)
	} else {
		fmt.Println("Invalid operation")
		return
	}
}
