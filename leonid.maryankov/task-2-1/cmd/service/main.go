package main

import "fmt"

func main() {
	var departments int
	_, err := fmt.Scan(&departments)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	for departments > 0 {
		var staff int
		_, err = fmt.Scan(&staff)
		if err != nil {
			fmt.Println("Error:", err)

			return
		}

		minTemp := 15
		maxTemp := 30

		for staff > 0 {
			var sign string
			var temp int

			_, err = fmt.Scan(&sign, &temp)
			if err != nil {
				fmt.Println("Error:", err)

				return
			}

			switch sign {
			case ">=":
				if temp > minTemp {
					minTemp = temp
				}
			case "<=":
				if temp < maxTemp {
					maxTemp = temp
				}
			default:
				fmt.Println("Error: invalid sign")

				return
			}

			staff--
		}

		if minTemp > maxTemp {
			fmt.Println(-1)
		} else {
			fmt.Println(minTemp)
		}

		departments--
	}
}
