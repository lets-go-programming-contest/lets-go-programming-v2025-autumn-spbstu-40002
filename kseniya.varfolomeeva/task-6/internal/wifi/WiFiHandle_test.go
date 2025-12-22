package wifi

import (
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var err error
	if args.Error(1) != nil {
		err = fmt.Errorf("mock error: %w", args.Error(1))
	}

	if args.Get(0) == nil {
		return nil, err
	}

	if ifaces, ok := args.Get(0).([]*wifi.Interface); ok {
		return ifaces, err
	}

	return nil, err
}
