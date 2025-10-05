package main

import "fmt"

func main() {
	var (
		departments int
		staff       int
		sign        string
		temp        int
		minTemp     int = 15
		maxTemp     int = 30
	)
	_, err := fmt.Scan(&departments)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for departments != 0 {
		_, err = fmt.Scan(&staff)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		for staff != 0 {
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
			if minTemp > maxTemp {
				fmt.Println(-1)
				return
			} else {
				fmt.Println(minTemp)
				staff--
			}
		}
		departments--
		minTemp, maxTemp = 15, 30
	}
}
