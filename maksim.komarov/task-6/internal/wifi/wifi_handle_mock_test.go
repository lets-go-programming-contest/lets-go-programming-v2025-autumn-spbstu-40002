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

	v := args.Get(0)
	if v != nil {
		typed, ok := v.([]*wifi.Interface)
		if ok {
			out = typed
		}
	}

	return out, args.Error(1)
}
