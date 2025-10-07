package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("5 + 3 ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	parts := strings.Split(input, " ")
	if len(parts) != 3 {
		fmt.Println("Invalid input format")
		return
	}

	a, err1 := strconv.ParseFloat(parts[0], 64)
	op := parts[1]
	b, err2 := strconv.ParseFloat(parts[2], 64)

	if err1 != nil {
		fmt.Println("Invalid first operand")
		return
	}
	if err2 != nil {
		fmt.Println("Invalid second operand")
		return
	}

	switch op {
	case "+":
		fmt.Printf(a+b)
	case "-":
		fmt.Printf(a-b)
	case "*":
		fmt.Printf(a*b)
	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Printf( a/b)
	default:
		fmt.Println("Invalid operation")
	}
}