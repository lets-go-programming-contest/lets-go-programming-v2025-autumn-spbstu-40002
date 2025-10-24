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

	firstOperandStr, _ := reader.ReadString('\n')
	firstOperandStr = strings.TrimSpace(firstOperandStr)

	secondOperandStr, _ := reader.ReadString('\n')
	secondOperandStr = strings.TrimSpace(secondOperandStr)

	operation, _ := reader.ReadString('\n')
	operation = strings.TrimSpace(operation)

	firstOperand, err := strconv.Atoi(firstOperandStr)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	secondOperand, err := strconv.Atoi(secondOperandStr)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	if operation != "+" && operation != "-" && operation != "*" && operation != "/" {
		fmt.Println("Invalid operation")
		return
	}

	var result int
	switch operation {
	case "+":
		result = firstOperand + secondOperand
	case "-":
		result = firstOperand - secondOperand
	case "*":
		result = firstOperand * secondOperand
	case "/":
		if secondOperand == 0 {
			fmt.Println("Division by zero")
			return
		}
		result = firstOperand / secondOperand
	}

	fmt.Println(result)
}
