package tempdata

import "errors"

var (
	ErrIncorrectRange = errors.New("minimal temperature must be greater maximum temperature")
	ErrInvalidSign    = errors.New("invalid sign")
)

type TemperatureData struct {
	optimalTemp int
	maxTemp     int
	minTemp     int
}

func NewTempData(optimalTemp, maxTemp, minTemp int) (*TemperatureData, error) {
	if minTemp > maxTemp {
		return nil, ErrIncorrectRange
	}

	return &TemperatureData{
		optimalTemp: optimalTemp,
		maxTemp:     maxTemp,
		minTemp:     minTemp,
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
