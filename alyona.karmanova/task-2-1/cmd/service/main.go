package main

import (
	"fmt"
)

func isCorrectInputTemperature(sign string, temp int) bool {
	return len(sign) == 2 && (sign[0] == '>' || sign[0] == '<') && sign[1] == '=' && (15 <= temp && temp <= 30)

}

func isCorrectInputCnt(c int) bool {
	return (c <= 1000 && c >= 1)
}

func applyLowerBound(currentMax *int, currentMin *int, newBoard int, currentTemp int) int {
	if newBoard > *currentMax {
		return -1
	}

	if newBoard > *currentMin {
		*currentMin = newBoard
		if currentTemp < newBoard {
			currentTemp = newBoard
		}
	}
	if currentTemp == -1 {
		currentTemp = newBoard
	}
	return currentTemp
}

func applyUpperBound(currentMax *int, currentMin *int, desiredTemp int, currentTemp int) int {
	if desiredTemp < *currentMin {
		return -1
	}

	if desiredTemp < *currentMax {
		*currentMax = desiredTemp
		if currentTemp > desiredTemp {
			currentTemp = desiredTemp
		}
	}
	if currentTemp == -1 {
		currentTemp = desiredTemp
	}
	return currentTemp

}
func main() {
	var countDepartment, countPeople int
	var sign string
	_, err := fmt.Scan(&countDepartment)
	if err != nil || !isCorrectInputCnt(countDepartment) {
		return
	}

	for i := 0; i < countDepartment; i++ {
		_, err = fmt.Scan(&countPeople)
		if err != nil || !isCorrectInputCnt(countPeople) {
			return
		}
		var tempMin, tempMax int = 14, 31
		var curTemp int = -1
		var newBoard int
		for j := 0; j < countPeople; j++ {
			_, err = fmt.Scan(&sign, &newBoard)
			if err != nil || !isCorrectInputTemperature(sign, newBoard) {
				return
			}
			isLookGreater := sign[0] == '>'
			if isLookGreater {
				curTemp = applyLowerBound(&tempMax, &tempMin, newBoard, curTemp)
			} else {
				curTemp = applyUpperBound(&tempMax, &tempMin, newBoard, curTemp)
			}
			fmt.Println(curTemp)
		}
	}
}
