package officeStruct

type Office struct {
	СurrentMax  int
	СurrentMin  int
	СurrentTemp int
}

func (o *Office) GetCurrentTemp() int {
	return o.СurrentTemp
}

func (o *Office) ApplyLowerBound(desiredTemp int) {
	if o.СurrentMin > o.СurrentMax {
		o.СurrentTemp = -1

		return
	}

	if desiredTemp > o.СurrentMax {
		o.СurrentMin = desiredTemp
		o.СurrentTemp = -1

		return
	}

	if desiredTemp > o.СurrentMin {
		o.СurrentMin = desiredTemp

		if o.СurrentTemp < desiredTemp {
			o.СurrentTemp = desiredTemp
		}
	}

	if o.СurrentTemp == -1 {
		o.СurrentTemp = -1
	}
}

func (o *Office) ApplyUpperBound(desiredTemp int) {
	if o.СurrentMin > o.СurrentMax {
		o.СurrentTemp = -1

		return
	}

	if desiredTemp < o.СurrentMin {
		o.СurrentMax = desiredTemp
		o.СurrentTemp = -1

		return
	}

	if desiredTemp < o.СurrentMax {
		o.СurrentMax = desiredTemp

		if o.СurrentTemp > desiredTemp {
			o.СurrentTemp = desiredTemp
		}
	}

	if o.СurrentTemp == -1 {
		o.СurrentTemp = -1
	}
}
