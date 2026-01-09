package main

import (
	"fmt"
)

const (
	minTempCond = 15
	maxTempCond = 30
	minDept     = 1
	maxDept     = 1000
	minEmp      = 1
	maxEmp      = 1000
)

type TempRange struct {
	min int
	max int
}

func newTempRange() *TempRange {
	return &TempRange{min: minTempCond, max: maxTempCond}
}

func (t *TempRange) update(op string, val int) bool {
	if val < minTempCond || val > maxTempCond {
		return false
	}
	switch op {
	case ">=":
		if val > t.min {
			t.min = val
		}
	case "<=":
		if val < t.max {
			t.max = val
		}
	default:
		return false
	}
	return true
}

func (t *TempRange) result() int {
	if t.min <= t.max {
		return t.min
	}
	return -1
}

func main() {
	var departments int
	_, err := fmt.Scan(&departments)
	if err != nil {
		return
	}

	if departments < minDept || departments > maxDept {
		return
	}

	for i := 0; i < departments; i++ {
		var employees int
		_, err := fmt.Scan(&employees)
		if err != nil {
			return
		}

		if employees < minEmp || employees > maxEmp {
			return
		}

		tr := newTempRange()

		for j := 0; j < employees; j++ {
			var op string
			var val int

			_, err := fmt.Scan(&op, &val)
			if err != nil {
				fmt.Println(-1)

				continue
			}

			ok := tr.update(op, val)
			if !ok {
				fmt.Println(-1)

			} else {
				fmt.Println(tr.result())

			}
		}
	}
}
