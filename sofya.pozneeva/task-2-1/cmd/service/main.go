package main

import "fmt"

func main() {
	var nSection uint16
	_, err := fmt.Scan(&nSection)
	if err != nil {
		fmt.Println("Invalid argument")
		
		return
	}
	
	for range nSection {
		var minTemp uint16 = 15
	
		var maxTemp uint16 = 30

		var nPeople uint16
		_, err = fmt.Scan(&nPeople)
		if err != nil {
			fmt.Println("Invalid argument")
			
			return
		}
		
		for range nPeople {
			var sign string
	
			var temp uint16
			
			_, err = fmt.Scanf("%s %d\n", &sign, &temp)
			if err != nil {
				fmt.Println("Invalid argument")
				
					return
			}	
			
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
