package main

import (
	"fmt"
)

func main() {
	var departments int
	_, err := fmt.Scan(&departments)
	if err != nil {
		return
	}

	for range departments {
		var minTemp = 15
		var maxTemp = 30

		var flag bool = false

		var workers int
		_, err = fmt.Scan(&workers)
		if err != nil {
			return
		}

		for i := 0; i < workers; i++ {
			var operator string
			_, err = fmt.Scan(&operator)
			if err != nil {
				return
			}

			var temp int
			_, err = fmt.Scan(&temp)
			if err != nil {
				return
			}

			switch operator {
			case ">=":
				if temp <= maxTemp && flag != true {
					minTemp = temp
					fmt.Println(minTemp)
				} else {
					flag = true
					fmt.Println(-1)
				}
			case "<=":
				if temp >= minTemp && flag != true {
					maxTemp = temp
					fmt.Println(minTemp)
				} else {
					flag = true
					fmt.Println(-1)
				}
			}
		}
	}
}
