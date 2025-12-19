package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/leonid.maryankov/task-6/internal/wifi"
	wifipkg "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errIface = errors.New("iface error")

type WiFiHandleMock struct {
	mock.Mock
}

func NewWiFiHandle(t *testing.T) *WiFiHandleMock {
	t.Helper()

	return &WiFiHandleMock{}
}

func (m *WiFiHandleMock) Interfaces() ([]*wifipkg.Interface, error) {
	args := m.Called()

	ifaces, _ := args.Get(0).([]*wifipkg.Interface)

	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("interfaces mock error: %w", err)
	}

	return ifaces, nil
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)

	ifaces := []*wifipkg.Interface{
		{
			Name:         "eth0",
			HardwareAddr: mustMAC("00:11:22:33:44:55"),
		},
		{
			Name:         "wlan0",
			HardwareAddr: mustMAC("aa:bb:cc:dd:ee:ff"),
		},
	}

	mockWiFi.On("Interfaces").Return(ifaces, nil)

	service := wifi.New(mockWiFi)

	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	require.Equal(t, []net.HardwareAddr{
		mustMAC("00:11:22:33:44:55"),
		mustMAC("aa:bb:cc:dd:ee:ff"),
	}, addrs)
}

func TestGetAddresses_Error(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)

	mockWiFi.On("Interfaces").Return(nil, errIface)

	service := wifi.New(mockWiFi)

	addrs, err := service.GetAddresses()

	require.Error(t, err)
	require.Nil(t, addrs)
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)

	ifaces := []*wifipkg.Interface{
		{Name: "eth0"},
		{Name: "wlan0"},
	}

	mockWiFi.On("Interfaces").Return(ifaces, nil)

	service := wifi.New(mockWiFi)

	names, err := service.GetNames()

	require.NoError(t, err)
	require.Equal(t, []string{"eth0", "wlan0"}, names)
}

func TestGetNames_Error(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)

	mockWiFi.On("Interfaces").Return(nil, errIface)

	service := wifi.New(mockWiFi)

	names, err := service.GetNames()

	require.Error(t, err)
	require.Nil(t, names)
}

func mustMAC(s string) net.HardwareAddr {
	m, err := net.ParseMAC(s)
	if err != nil {
		panic(err)
	}

	return m
}
