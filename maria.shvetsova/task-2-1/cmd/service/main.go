package main

import (
	"fmt"
)

func main() {
	var n int
	fmt.Scan(&n)
	if n < 1 || n > 1000 {
		return
	}

	for i := 0; i < n; i++ {
		var k int
		fmt.Scan(&k)
		if k < 1 || k > 1000 {
			return
		}

		var sign string
		var t int
		maxTemp := 30
		minTemp := 15
		optimalTemp := 15
		for j := 0; j < k; j++ {
			fmt.Scan(&sign)
			fmt.Scan(&t)
			if sign == ">=" {
				if t <= maxTemp {
					if t <= minTemp || t <= optimalTemp {
						minTemp = t
					} else {
						optimalTemp = t
					}
					fmt.Println(optimalTemp)
				} else {
					minTemp = t
					fmt.Println(-1)
				}
			} else {
				if t >= minTemp {
					if t >= maxTemp || t >= optimalTemp {
						maxTemp = t
					} else {
						optimalTemp = t
					}
					fmt.Println(optimalTemp)
				} else {
					maxTemp = t
					fmt.Println(-1)
				}
			}
		}
	}
}
