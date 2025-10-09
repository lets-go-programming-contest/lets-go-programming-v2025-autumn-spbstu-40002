package tempdata

import "errors"

var ErrInvalidSign = errors.New("invalid sign")

type TemperatureData struct {
	optimalTemp int
	maxTemp     int
	minTemp     int
}

func NewTempData() (*TemperatureData, error) {
	return &TemperatureData{
		optimalTemp: 15,
		maxTemp:     30,
		minTemp:     15,
	}, nil
}

func (td *TemperatureData) GetOptimalTemp() int {
	return td.optimalTemp
}

func (td *TemperatureData) ChangeOptimalTemp(sign string, temp int) error {
	switch sign {
	case ">=":
		if temp > td.maxTemp {
			td.optimalTemp = -1
		} else if temp > td.minTemp {
			td.minTemp = temp
		}
	case "<=":
		if temp < td.minTemp {
			td.optimalTemp = -1
		} else if temp < td.maxTemp {
			td.maxTemp = temp
		}
	default:
		return ErrInvalidSign
	}

	if td.optimalTemp != -1 {
		td.optimalTemp = td.minTemp
	}

	return nil
}
