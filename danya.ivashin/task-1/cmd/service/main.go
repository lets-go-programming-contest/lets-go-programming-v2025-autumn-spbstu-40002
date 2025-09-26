package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		fmt.Println("Invalid first operand")
		return
	}
	firstStr := strings.TrimSpace(scanner.Text())
	first, err := strconv.Atoi(firstStr)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	if !scanner.Scan() {
		fmt.Println("Invalid second operand")
		return
	}
	secondStr := strings.TrimSpace(scanner.Text())
	second, err := strconv.Atoi(secondStr)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	if !scanner.Scan() {
		fmt.Println("Invalid operation")
		return
	}
	op := strings.TrimSpace(scanner.Text())

	switch op {
	case "+":
		fmt.Println(first + second)
	case "-":
		fmt.Println(first - second)
	case "*":
		fmt.Println(first * second)
	case "/":
		if second == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(first / second)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
