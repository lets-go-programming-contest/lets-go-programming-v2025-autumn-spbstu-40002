package main

import (
	"errors"
	"fmt"
)

type Conditioner struct {
	MinTemp uint16
	MaxTemp uint16
}

var errInvalidTemp = errors.New("invalid temperature")

func (cond *Conditioner) GetTemp() error {
	if cond.MinTemp <= cond.MaxTemp {
		return nil
	}

	return errInvalidTemp
}

var errInvalidSign = errors.New("invalid sign")

func (cond *Conditioner) SetTemp(sign string, temp uint16) error {
	switch sign {
	case ">=":
		if cond.MinTemp < temp {
			cond.MinTemp = temp
		}
	case "<=":
		if cond.MaxTemp > temp {
			cond.MaxTemp = temp
		}
	default:
		return errInvalidSign
	}

	return nil
}

var errInvalidArgument = errors.New("invalid argument")

func main() {
	var nSection uint16

	_, err := fmt.Scan(&nSection)
	if err != nil {
		fmt.Println(errInvalidArgument)

		return
	}

	const minCond uint16 = 15

	const maxCond uint16 = 30

	for range nSection {
		var nPeople uint16

		_, err = fmt.Scan(&nPeople)
		if err != nil {
			fmt.Println(errInvalidArgument)

			return
		}

		pCond := &Conditioner{MinTemp: minCond, MaxTemp: maxCond}

		for range nPeople {
			var sign string

			var temp uint16

			_, err = fmt.Scan(&sign)
			if err != nil {
				fmt.Println(errInvalidArgument)

				return
			}

			_, err = fmt.Scan(&temp)
			if err != nil {
				fmt.Println(errInvalidArgument)

				return
			}

			err = pCond.SetTemp(sign, temp)
			if err != nil {
				fmt.Println(err)

				return
			}

			err = pCond.GetTemp()
			if err != nil {
				fmt.Println(-1)
			} else {
				fmt.Println(pCond.MinTemp)
			}
		}
	}
}
