package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"

	"github.com/Svyatoslav2324/task-6/internal/mocks"
	myWifi "github.com/Svyatoslav2324/task-6/internal/wifi"
)

func TestGetAddresses_Success(t *testing.T) {
	mockWiFi := mocks.NewWiFiHandle(t)

	ifaces := []*wifi.Interface{
		{
			Name:         "wlan0",
			HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		},
		{
			Name:         "wlan1",
			HardwareAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		},
	}

	mockWiFi.
		On("Interfaces").
		Return(ifaces, nil).
		Once()

	service := myWifi.New(mockWiFi)

	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	require.Equal(t, []net.HardwareAddr{
		ifaces[0].HardwareAddr,
		ifaces[1].HardwareAddr,
	}, addrs)
}

func TestGetAddresses_Error(t *testing.T) {
	mockWiFi := mocks.NewWiFiHandle(t)

	mockWiFi.
		On("Interfaces").
		Return([]*wifi.Interface(nil), errors.New("fail")).
		Once()

	service := myWifi.New(mockWiFi)

	addrs, err := service.GetAddresses()

	require.Error(t, err)
	require.Nil(t, addrs)
	require.Contains(t, err.Error(), "getting interfaces")
}

func TestGetNames_Success(t *testing.T) {
	mockWiFi := mocks.NewWiFiHandle(t)

	ifaces := []*wifi.Interface{
		{Name: "wlan0"},
		{Name: "wlan1"},
	}

	mockWiFi.
		On("Interfaces").
		Return(ifaces, nil).
		Once()

	service := myWifi.New(mockWiFi)

	names, err := service.GetNames()

	require.NoError(t, err)
	require.Equal(t, []string{"wlan0", "wlan1"}, names)
}

func TestGetNames_Error(t *testing.T) {
	mockWiFi := mocks.NewWiFiHandle(t)

	mockWiFi.
		On("Interfaces").
		Return([]*wifi.Interface(nil), errors.New("fail")).
		Once()

	service := myWifi.New(mockWiFi)

	names, err := service.GetNames()

	require.Error(t, err)
	require.Nil(t, names)
	require.Contains(t, err.Error(), "getting interfaces")
}
