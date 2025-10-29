package main

import (
	"errors"
	"fmt"
)

const (
	MinTemp = 15
	MaxTemp = 30
)

var errSign = errors.New("invalid sign")

type Temperature struct {
	minTemp int
	maxTemp int
}

func (t *Temperature) changeTemp(sign string, temp int) error {
	switch sign {
	case ">=":
		if temp > t.minTemp {
			t.minTemp = temp
		}
	case "<=":
		if temp < t.maxTemp {
			t.maxTemp = temp
		}
	default:
		return errSign
	}

	return nil
}

func (t *Temperature) getTemp() int {
	if t.minTemp > t.maxTemp {
		return -1
	}

	return t.minTemp
}

func main() {
	var (
		departments int
		staff       int
		sign        string
		temp        int
	)

	_, err := fmt.Scan(&departments)
	if err != nil {
		fmt.Println("", err)

		return
	}

	for departments > 0 {
		_, err = fmt.Scan(&staff)
		if err != nil {
			fmt.Println("", err)

			return
		}

		dept := Temperature{
			minTemp: MinTemp,
			maxTemp: MaxTemp,
		}

		for staff > 0 {
			_, err = fmt.Scan(&sign, &temp)
			if err != nil {
				fmt.Println("", err)

				return
			}

			err = dept.changeTemp(sign, temp)
			if err != nil {
				fmt.Println(-1)
			} else {
				fmt.Println(dept.getTemp())
			}

			staff--
		}

		departments--
	}
}
