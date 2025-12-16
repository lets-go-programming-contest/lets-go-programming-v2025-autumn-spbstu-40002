package wifi

import (
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"

	"task-6/internal/wifi/mocks"
)

func TestWiFiService_GetAddresses_OK(t *testing.T) {
	m := mocks.NewWiFiHandle(t)

	ifaces := []*wifi.Interface{
		{Name: "wlan0", HardwareAddr: net.HardwareAddr{0, 1, 2, 3, 4, 5}},
		{Name: "wlan1", HardwareAddr: net.HardwareAddr{10, 11, 12, 13, 14, 15}},
	}

	m.On("Interfaces").Return(ifaces, nil).Once()

	service := New(m)
	got, err := service.GetAddresses()
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, ifaces[0].HardwareAddr, got[0])
	require.Equal(t, ifaces[1].HardwareAddr, got[1])
}

func TestWiFiService_GetAddresses_Error(t *testing.T) {
	m := mocks.NewWiFiHandle(t)

	m.On("Interfaces").Return(nil, errors.New("boom")).Once()

	service := New(m)
	_, err := service.GetAddresses()
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "getting interfaces:"))
}

func TestWiFiService_GetNames_OK(t *testing.T) {
	m := mocks.NewWiFiHandle(t)

	ifaces := []*wifi.Interface{
		{Name: "wlan0", HardwareAddr: net.HardwareAddr{0, 1, 2, 3, 4, 5}},
		{Name: "wlan1", HardwareAddr: net.HardwareAddr{10, 11, 12, 13, 14, 15}},
	}

	m.On("Interfaces").Return(ifaces, nil).Once()

	service := New(m)
	got, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"wlan0", "wlan1"}, got)
}

func TestWiFiService_GetNames_Error(t *testing.T) {
	m := mocks.NewWiFiHandle(t)

	m.On("Interfaces").Return(nil, errors.New("boom")).Once()

	service := New(m)
	_, err := service.GetNames()
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "getting interfaces:"))
}
