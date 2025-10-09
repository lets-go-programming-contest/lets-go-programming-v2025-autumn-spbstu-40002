package main

import "fmt"

func main() {
	var minTemp uint8
	
	var maxTemp uint8
	
	var nSection uint8
	
	var nPeople uint8
	
	var sign string
	
	var temp uint8
	
	var err error
	
	_, err = fmt.Scanln(&nSection)
	
	if err != nil {
		fmt.Print("Invalid argument")
		
		return
	}
	
	for range nSection {
		minTemp = 15
		maxTemp = 30
		_, err = fmt.Scanln(&nPeople)
			
		if err != nil {
			fmt.Print("Invalid argument")
			
			return
		}
		
		for range nPeople {
			_, err = fmt.Scanf("%s %d\n", &sign, &temp)
					
			if err != nil {
				fmt.Print("Invalid argument")
				
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
