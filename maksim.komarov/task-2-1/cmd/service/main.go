package main

import "fmt"

func main() {
	const minT = 15
	const maxT = 30

	var n, k int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	if _, err := fmt.Scan(&k); err != nil {
		return
	}

	for dept := 0; dept < n; dept++ {
		low, high := minT, maxT

		for i := 0; i < k; i++ {
			var op string
			var v int
			if _, err := fmt.Scan(&op, &v); err != nil {
				return
			}

			switch op {
			case ">=":
				if v > low {
					low = v
				}
			case "<=":
				if v < high {
					high = v
				}
			}

			if low > high {
				fmt.Println(-1)
			} else {
				fmt.Println(low)
			}
		}
	}
}
