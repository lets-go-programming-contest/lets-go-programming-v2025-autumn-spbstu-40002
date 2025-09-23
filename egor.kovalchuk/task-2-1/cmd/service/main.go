package main

import "fmt"

func main() {
	var N int
	if _, err := fmt.Scan(&N); err != nil {
		return
	}

	for i := 0; i < N; i++ {
		min := 15
		max := 30
		var K int

		if _, err := fmt.Scan(&K); err != nil {
			return
		}

		contradiction := false

		for j := 0; j < K; j++ {
			var operator string
			var temperature int

			if _, err := fmt.Scan(&operator, &temperature); err != nil {
				return
			}

			if contradiction {
				fmt.Println(-1)
				continue
			}

			switch operator {
			case ">=":
				if temperature > min {
					min = temperature
				}
			case "<=":
				if temperature < max {
					max = temperature
				}
			}

			if min > max {
				contradiction = true
				fmt.Println(-1)
			} else {
				fmt.Println(min)
			}
		}
	}
}
