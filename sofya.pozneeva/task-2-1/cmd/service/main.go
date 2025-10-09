package main

import "fmt"

var {
	minTemp uint8
	maxTemp uint8
	nSection uint8
	nPeople uint8
	sign string
	temp uint8
	err error
}

func main() {
	_, err = fmt.Scan(&nSection)
	
	if err != nil {
		fmt.Print("Invalid argument")
		return
	} else {
		for range nSection {
			minTemp = 15
			maxTemp = 30
			_, err = fmt.Scan(&nPeople)
			
			if err != nil {
				fmt.Print("Invalid argument")
				return
			} else {
				for range nPeople {
					_, err = fmt.Scanf("\n%s %d", &sign, &temp)
					
					if err != nil {
						fmt.Print("Invalid argument")
						return
					} else {
						
						if sign == ">=" {
							
							if minTemp < temp {
								minTemp = temp
							}
						} else {
							
							if maxTemp > temp {
								maxTemp = temp
							}
						}
						
						if minTemp <= maxTemp {
							fmt.Println(minTemp)
						} else {
							fmt.Println(-1)
						}
					}
				}
			}
		}
	}
}
