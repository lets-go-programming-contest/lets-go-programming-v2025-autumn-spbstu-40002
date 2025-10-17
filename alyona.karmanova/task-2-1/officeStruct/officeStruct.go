package officestruct

type Office struct {
	CurrentMax  int
	CurrentMin  int
	CurrentTemp int
}

func (o *Office) GetCurrentTemp() int {
	return o.CurrentTemp
}

func (o *Office) ApplyLowerBound(desiredTemp int) {
	if o.CurrentMin > o.CurrentMax {
		o.CurrentTemp = -1

		return
	}

	if desiredTemp > o.CurrentMax {
		o.CurrentMin = desiredTemp
		o.CurrentTemp = -1

		return
	}

	if desiredTemp > o.CurrentMin {
		o.CurrentMin = desiredTemp

		if o.CurrentTemp < desiredTemp {
			o.CurrentTemp = desiredTemp
		}
	}

	if o.CurrentTemp == -1 {
		o.CurrentTemp = -1
	}
}

func (o *Office) ApplyUpperBound(desiredTemp int) {
	if o.CurrentMin > o.CurrentMax {
		o.CurrentTemp = -1

		return
	}

	if desiredTemp < o.CurrentMin {
		o.CurrentMax = desiredTemp
		o.CurrentTemp = -1

		return
	}

	if desiredTemp < o.CurrentMax {
		o.CurrentMax = desiredTemp

		if o.CurrentTemp > desiredTemp {
			o.CurrentTemp = desiredTemp
		}
	}

	if o.CurrentTemp == -1 {
		o.CurrentTemp = -1
	}
}
