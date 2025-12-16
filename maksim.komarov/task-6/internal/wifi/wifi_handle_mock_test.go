package wifi_test

import (
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type wiFiHandleMock struct {
	mock.Mock
}

func (m *wiFiHandleMock) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var out []*wifi.Interface

	v0 := args.Get(0)
	if v0 != nil {
		typed, ok := v0.([]*wifi.Interface)
		if ok {
			out = typed
		}
	}

	var err error

	v1 := args.Get(1)
	if v1 != nil {
		typedErr, ok := v1.(error)
		if ok {
			err = typedErr
		}
	}

	return out, err
}
