package main

import (
	"fmt"
)

const (
	maxDep = 1000
	minDep = 1
)

func main() {
	var Departaments, employees int
	_, err := fmt.Scan(&Departaments)
	if err != nil || Departaments > maxDep || Departaments < minDep {
		fmt.Println("Invalid input")
		return
	}

	for i := 1; i < Departaments+1; i++ {
		var sign string
		var personalTemperature int
		minTemp, maxTemp := 15, 30
		_, err = fmt.Scan(&employees)
		if err != nil || employees > maxDep || employees < minDep {
			fmt.Println("Invalid input")
			return
		}
		for j := 1; j < employees+1; j++ {
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
				if personalTemperature >= minTemp {
					minTemp = personalTemperature
				}
				if personalTemperature <= minTemp {
					fmt.Println(minTemp)
					continue
				}
				fmt.Println(minTemp)
			case "<=":
				if personalTemperature < minTemp {
					fmt.Println("-1")
					return
				}
				if personalTemperature >= maxTemp {
					fmt.Println(maxTemp)
					continue
				}
				if personalTemperature < maxTemp {
					maxTemp = personalTemperature
				}
				fmt.Println(minTemp)
			}
		}
	}

}
