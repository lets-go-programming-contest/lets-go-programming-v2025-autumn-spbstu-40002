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
	n, _ := strconv.Atoi(scanner.Text())
	scanner.Scan()
	employeeCount, _ := strconv.Atoi(scanner.Text())

	for department := 0; department < n; department++ {
		minTemp := 15
		maxTemp := 30
		possible := true

		for employee := 0; employee < employeeCount; employee++ {
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

			if minTemp <= maxTemp && possible {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
				possible = false
			}
		}
	}
}
