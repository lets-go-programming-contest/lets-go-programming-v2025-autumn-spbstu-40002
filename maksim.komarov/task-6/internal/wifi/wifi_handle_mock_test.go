package wifi

import (
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type WiFiHandleMock struct {
	mock.Mock
}

func (m *WiFiHandleMock) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	var r0 []*wifi.Interface
	if v := args.Get(0); v != nil {
		r0 = v.([]*wifi.Interface)
	}
	return r0, args.Error(1)
}
