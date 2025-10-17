package main

import (
	"fmt"
)

const (
	maxDep = 1000
	minDep = 1
)

func main() {
	var departaments, employees int
	
	_, err := fmt.Scan(&departaments)
	if err != nil || departaments > maxDep || departaments < minDep {
		fmt.Println("Invalid input")
		return
	}

	for i := 0; i < departaments; i++ {
		var sign string
		var personalTemperature int
		minTemp, maxTemp := 15, 30
		
		_, err = fmt.Scan(&employees)
		if err != nil || employees > maxDep || employees < minDep {
			fmt.Println("Invalid input")
			return
		}

		for j := 0; j < employees; j++ {
			_, err = fmt.Scan(&sign, &personalTemperature)
			if err != nil {
				fmt.Println("Invalid input")
				return
			}

			switch sign {
			case ">=":
				if personalTemperature > maxTemp {
					fmt.Println("-1")
					return
				}
				
				if personalTemperature > minTemp {
					minTemp = personalTemperature
				}
				
				fmt.Println(minTemp)

			case "<=":
				if personalTemperature < minTemp {
					fmt.Println("-1")
					return
				}
				
				if personalTemperature < maxTemp {
					maxTemp = personalTemperature
				}
				
				fmt.Println(minTemp)
			}
		}
	}
}