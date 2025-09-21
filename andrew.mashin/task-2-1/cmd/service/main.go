package main

import (
	"fmt"
)

func main() {
	var n int
	var k int
	var operator string
	var temp int
	var minTemp int
	var maxTemp int

	_, err := fmt.Scan(&n)
	if err != nil {
		return
	}

	for i := 0; i < n; i++ {
		minTemp = 15
		maxTemp = 30

		_, err = fmt.Scan(&k)
		if err != nil {
			return
		}

		for j := 0; j < k; j++ {
			_, err = fmt.Scan(&operator)
			if err != nil {
				return
			}

			_, err = fmt.Scan(&temp)
			if err != nil {
				return
			}

			if operator == ">=" && temp <= maxTemp {
				minTemp = temp
				fmt.Println(minTemp)
			} else if operator == "<=" && temp >= minTemp {
				maxTemp = temp
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
