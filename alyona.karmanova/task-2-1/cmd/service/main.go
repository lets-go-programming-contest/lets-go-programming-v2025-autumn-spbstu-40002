package main

import (
	"fmt"

	officestruct "github.com/HuaChenju/task-2-1/officeStruct"
)

const (
	MinTempConst   = 15
	MaxTempConst   = 30
	MinCorrectData = 1
	MaxCorrectData = 1000
)

func isDateCorrect(c int) bool {
	return (c <= MaxCorrectData && c >= MinCorrectData)
}

func main() {
	var (
		countDepartment, countPeople, newBoard int
		sign                                   string
	)

	_, err := fmt.Scan(&countDepartment)

	if err != nil || !isDateCorrect(countDepartment) {
		fmt.Println("invalid number of departments")

		return
	}

	for range countDepartment {
		_, err = fmt.Scan(&countPeople)

		if err != nil || !isDateCorrect(countPeople) {
			fmt.Println("invalid number of employees")

			return
		}

		room := officestruct.Office{
			CurrentMax:  MaxTempConst,
			CurrentMin:  MinTempConst,
			CurrentTemp: MinTempConst,
		}

		for range countPeople {
			_, err = fmt.Scan(&sign, &newBoard)
			if err != nil {
				return
			}

			isLookGreater := sign[0] == '>'

			if isLookGreater {
				room.ApplyLowerBound(newBoard)
			} else {
				room.ApplyUpperBound(newBoard)
			}

			fmt.Println(room.GetCurrentTemp())
		}
	}
}
