package main

import "fmt"

const (
	defaultLow  = 15
	defaultHigh = 30
)

func main() {
	var departments int
	if _, err := fmt.Scan(&departments); err != nil {
		return
	}
	if departments < 0 {
		return
	}

	depts := make([]struct{}, departments)
	for range depts {
		low := defaultLow
		high := defaultHigh

		var employeeCount int
		if _, err := fmt.Scan(&employeeCount); err != nil {
			return
		}
		if employeeCount < 0 {
			return
		}

		contradiction := false

		emps := make([]struct{}, employeeCount)
		for range emps {
			operator := ""
			temperature := 0

			_, err := fmt.Scan(&operator, &temperature)
			if err != nil {
				return
			}

			if contradiction {
				fmt.Println(-1)

				continue
			}

			if operator != ">=" && operator != "<=" {
				contradiction = true
				fmt.Println(-1)

				continue
			}

			switch operator {
			case ">=":
				if temperature > low {
					low = temperature
				}
			case "<=":
				if temperature < high {
					high = temperature
				}
			}

			if low > high {
				contradiction = true

				fmt.Println(-1)

				continue
			}

			fmt.Println(low)
		}
	}
}
