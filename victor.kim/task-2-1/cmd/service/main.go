package main

import (
	"fmt"
	"os"
)

func main() {
	var deptNum int
	if _, err := fmt.Scan(&deptNum); err != nil {
		return
	}

	for range deptNum {
		var emplNum int
		if _, err := fmt.Scan(&emplNum); err != nil {
			fmt.Fprintln(os.Stderr, "couldn't read the department count:", err)
		}

		lowerBound, higherBound := 15, 30

		for range emplNum {
			var (
				operator    string
				desiredTemp int
			)

			if _, err := fmt.Scan(&operator, &desiredTemp); err != nil {
				fmt.Fprintln(os.Stderr, "couldn't read employee:", err)

				return
			}

			switch operator {
			case "<=":
				if desiredTemp < higherBound {
					higherBound = desiredTemp
				}
			case ">=":
				if desiredTemp > lowerBound {
					lowerBound = desiredTemp
				}
			default:
				fmt.Fprintln(os.Stderr, "unsupported operator:", operator)

				return
			}

			if lowerBound > higherBound {
				fmt.Println(-1)
			} else {
				fmt.Println(lowerBound)
			}
		}
	}
}
