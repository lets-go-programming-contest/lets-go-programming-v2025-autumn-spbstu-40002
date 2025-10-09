package main

import (
	"fmt"
)

const (
	MinTempConst   = 15
	MaxTempConst   = 30
	MinCorrectData = 1
	MaxCorrectData = 1000
)

type Office struct {
	currentMax  int
	currentMin  int
	currentTemp int
}

func (o *Office) GetCurrentTemp() int {
	return o.currentTemp
}

func (o *Office) applyLowerBound(desiredTemp int) {
	if o.currentMin > o.currentMax {
		o.currentTemp = -1
		return
	}
	if desiredTemp > o.currentMax {
		o.currentTemp = -1
		return
	}

	if desiredTemp > o.currentMin {
		o.currentMin = desiredTemp

		if o.currentTemp < desiredTemp {
			o.currentTemp = desiredTemp
		}
	}

	if o.currentTemp == -1 {
		o.currentTemp = desiredTemp
	}
}

func isDateCorrect(c int) bool {
	return (c <= MaxCorrectData && c >= MinCorrectData)
}

func (o *Office) applyUpperBound(desiredTemp int) {
	if o.currentMin > o.currentMax {
		o.currentMin = -1
		return
	}
	if desiredTemp < o.currentMin {
		o.currentMax = desiredTemp
		o.currentMin = -1
		return
	}

	if desiredTemp < o.currentMax {
		o.currentMax = desiredTemp

		if o.currentTemp > desiredTemp {
			o.currentTemp = desiredTemp
		}
	}

	if o.currentTemp == -1 {
		o.currentTemp = desiredTemp
	}
}

func main() {
	var (
		countDepartment, countPeople, newBoard int
		sign                                   string
	)

	room := Office{
		currentMax:  30,
		currentMin:  15,
		currentTemp: 15,
	}

	_, err := fmt.Scan(&countDepartment)

	if err != nil || !isDateCorrect(countDepartment) {
		return
	}

	for range countDepartment {
		_, err = fmt.Scan(&countPeople)

		if err != nil || !isDateCorrect(countPeople) {
			return
		}

		for range countPeople {
			_, err = fmt.Scan(&sign, &newBoard)

			if err != nil {
				return
			}

			isLookGreater := sign[0] == '>'

			if isLookGreater {
				room.applyLowerBound(newBoard)
			} else {
				room.applyUpperBound(newBoard)
			}

			fmt.Println(room.GetCurrentTemp())
		}
	}
}
