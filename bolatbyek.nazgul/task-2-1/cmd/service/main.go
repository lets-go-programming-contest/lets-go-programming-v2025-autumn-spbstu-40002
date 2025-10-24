package main

import (
	"fmt"
	"strings"
)

func main() {
	var N, K int
	fmt.Scan(&N, &K)

	for i := 0; i < N; i++ {
		// Process each department
		minTemp := 15
		maxTemp := 30
		
		for j := 0; j < K; j++ {
			var constraint string
			fmt.Scan(&constraint)
			
			// Parse constraint (e.g., ">=30" or "<=25")
			if strings.HasPrefix(constraint, ">=") {
				var temp int
				fmt.Sscanf(constraint, ">=%d", &temp)
				if temp > minTemp {
					minTemp = temp
				}
			} else if strings.HasPrefix(constraint, "<=") {
				var temp int
				fmt.Sscanf(constraint, "<=%d", &temp)
				if temp < maxTemp {
					maxTemp = temp
				}
			}
			
			// Check if valid range exists
			if minTemp > maxTemp {
				fmt.Println(-1)
			} else {
				// Output the minimum valid temperature
				fmt.Println(minTemp)
			}
		}
	}
}