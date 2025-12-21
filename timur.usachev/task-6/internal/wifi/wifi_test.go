package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	wifipkg "github.com/t1wt/task-6/internal/wifi"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	errIfaces           = errors.New("getting interfaces failed")
	errInvalidIfaceType = errors.New("invalid interface type")
)

type MockWiFi struct {
	mock.Mock
}

func (m *MockWiFi) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	ifaceSlice, ok := args.Get(0).([]*wifi.Interface)

	if !ok {
		return nil, errInvalidIfaceType
	}

	if err := args.Error(1); err != nil {
		return ifaceSlice, fmt.Errorf("%w", err)
	}

	return ifaceSlice, nil
}

func TestGetAddresses_Success(t *testing.T) {
	t.Parallel()

	hw := net.HardwareAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	if1 := &wifi.Interface{Name: "wlan0", HardwareAddr: hw}

	m := new(MockWiFi)
	m.On("Interfaces").Return([]*wifi.Interface{if1}, nil)

	service := wifipkg.New(m)
	addrs, err := service.GetAddresses()
	require.NoError(t, err)
	require.Equal(t, []net.HardwareAddr{hw}, addrs)
	m.AssertExpectations(t)
}

func TestGetAddresses_Empty(t *testing.T) {
	t.Parallel()

	m := new(MockWiFi)
	m.On("Interfaces").Return([]*wifi.Interface{}, nil)

	service := wifipkg.New(m)
	addrs, err := service.GetAddresses()
	require.NoError(t, err)
	require.Empty(t, addrs)
	m.AssertExpectations(t)
}

func TestGetAddresses_Error(t *testing.T) {
	t.Parallel()

	m := new(MockWiFi)
	m.On("Interfaces").Return(([]*wifi.Interface)(nil), errIfaces)

	service := wifipkg.New(m)
	addrs, err := service.GetAddresses()
	require.Nil(t, addrs)
	require.ErrorIs(t, err, errIfaces)
	m.AssertExpectations(t)
}

func TestGetNames_Success(t *testing.T) {
	t.Parallel()

	if1 := &wifi.Interface{Name: "wlan0"}
	if2 := &wifi.Interface{Name: "eth0"}

	m := new(MockWiFi)
	m.On("Interfaces").Return([]*wifi.Interface{if1, if2}, nil)

	service := wifipkg.New(m)
	names, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"wlan0", "eth0"}, names)
	m.AssertExpectations(t)
}

func TestGetNames_Empty(t *testing.T) {
	t.Parallel()

	m := new(MockWiFi)
	m.On("Interfaces").Return([]*wifi.Interface{}, nil)

	service := wifipkg.New(m)
	names, err := service.GetNames()
	require.NoError(t, err)
	require.Empty(t, names)
	m.AssertExpectations(t)
}

func TestGetNames_Error(t *testing.T) {
	t.Parallel()

	m := new(MockWiFi)
	m.On("Interfaces").Return(([]*wifi.Interface)(nil), errIfaces)

	service := wifipkg.New(m)
	names, err := service.GetNames()
	require.Nil(t, names)
	require.ErrorIs(t, err, errIfaces)
	m.AssertExpectations(t)
}
