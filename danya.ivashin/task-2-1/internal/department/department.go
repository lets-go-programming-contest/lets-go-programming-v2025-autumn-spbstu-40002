package department

import "errors"

var errInvalidOperation = errors.New("invalid operator")

const (
	UpperLimit = 30
	LowerLimit = 15
)

type ComparisonOperator string

const (
	LessEqual    ComparisonOperator = "<="
	GreaterEqual ComparisonOperator = ">="
)

type ChangeRequest struct {
	Operator    ComparisonOperator
	Temperature int
}

type Department struct {
	minT, maxT int
}

func New() *Department {
	return &Department{
		minT: LowerLimit,
		maxT: UpperLimit,
	}
}

func (d *Department) setMax(temp int) {
	if d.maxT > temp {
		d.maxT = temp
	}
}

func (d *Department) setMin(temp int) {
	if d.minT < temp {
		d.minT = temp
	}
}

func (d *Department) Recalculate(req *ChangeRequest) (int, error) {
	switch req.Operator {
	case LessEqual:
		d.setMax(req.Temperature)
	case GreaterEqual:
		d.setMin(req.Temperature)
	default:
		return -1, errInvalidOperation
	}

	return d.Optimum(), nil
}

func (d *Department) Optimum() int {
	if d.maxT < d.minT {
		return -1
	}

	return d.minT
}
