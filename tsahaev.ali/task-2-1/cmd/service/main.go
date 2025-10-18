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

	scanner.Scan()
	departmentsCount, _ := strconv.Atoi(scanner.Text())

	scanner.Scan()
	employeesCount, _ := strconv.Atoi(scanner.Text())

	for range departmentsCount {
		minTemp := 15
		maxTemp := 30
		results := make([]int, employeesCount)

		for i := range employeesCount {
			scanner.Scan()
			line := scanner.Text()
			parts := strings.Split(line, " ")

			operator := parts[0]
			temp, _ := strconv.Atoi(parts[1])

			if operator == ">=" {
				if temp > minTemp {
					minTemp = temp
				}
			} else if operator == "<=" {
				if temp < maxTemp {
					maxTemp = temp
				}
			}

			if minTemp <= maxTemp {
				results[i] = minTemp
			} else {
				results[i] = -1
			}
		}

		for _, result := range results {
			fmt.Println(result)
		}

		if departmentsCount > 1 {
			fmt.Println()
		}
	}
}
