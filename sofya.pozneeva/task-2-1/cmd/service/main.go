package main

import "fmt"

func main() {
	
	var minTemp uint16
	
	var maxTemp uint16
	
	var nSection uint16
	
	var nPeople uint16
	
	var sign string
	
	var temp uint16
	
	var err error
	
	_, err = fmt.Scanln(&nSection)
	
	if err != nil {
		fmt.Println("Invalid argument")
		
		return
	}
	
	for range nSection {
		minTemp = 15
		maxTemp = 30
		_, err = fmt.Scanln(&nPeople)
			
		if err != nil {
			fmt.Println("Invalid argument")
			
			return
		}
		
		for range nPeople {
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
