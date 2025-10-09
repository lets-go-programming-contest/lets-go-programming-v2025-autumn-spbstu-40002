package main

import "fmt"

func main() {
	var min uint8
	var max uint8
	var nSection uint8
	var nPeople uint8
	var sign string
	var temp uint8
	fmt.Scan(&nSection)
	for i := 0; i < int(nSection); i++ {
		min = 15
		max = 30
		fmt.Scan(&nPeople)
		for j := 0; j < int(nPeople); j++ {
			fmt.Scanf("\n%s %d", &sign, &temp)
			if sign == ">=" {
				if min < temp {
					min = temp
				}
			} else {
				if max > temp {
					max = temp
				}
			}
			if min < max {
				fmt.Println(min)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
