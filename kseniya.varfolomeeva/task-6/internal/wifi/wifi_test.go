package wifi_test

import (
	"errors"
	"net"
	"testing"

	"task-6/internal/wifi"
	mdlayherwifi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

var ErrExpected = errors.New("expected error")

func TestNew(t *testing.T) {
	t.Parallel()

	mockWiFi := new(wifi.MockWiFiHandle)
	service := wifi.New(mockWiFi)

	require.Equal(t, mockWiFi, service.WiFi, "Expected WiFi to be set")
	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetAddresses_Success(t *testing.T) {
	t.Parallel()

	mockWiFi := new(wifi.MockWiFiHandle)

	hwAddr1, err := net.ParseMAC("00:11:22:33:44:55")
	require.NoError(t, err)
	hwAddr2, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	require.NoError(t, err)

	interfaces := []*mdlayherwifi.Interface{
		{Name: "wlan0", HardwareAddr: hwAddr1},
		{Name: "wlan1", HardwareAddr: hwAddr2},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil).Once()

	service := wifi.New(mockWiFi)
	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	require.Len(t, addrs, 2)
	require.Equal(t, hwAddr1, addrs[0])
	require.Equal(t, hwAddr2, addrs[1])

	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetAddresses_Empty(t *testing.T) {
	t.Parallel()

	mockWiFi := new(wifi.MockWiFiHandle)
	interfaces := []*mdlayherwifi.Interface{}

	mockWiFi.On("Interfaces").Return(interfaces, nil).Once()

	service := wifi.New(mockWiFi)
	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	require.Empty(t, addrs)

	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetAddresses_Error(t *testing.T) {
	t.Parallel()

	mockWiFi := new(wifi.MockWiFiHandle)
	mockWiFi.On("Interfaces").Return([]*mdlayherwifi.Interface{}, ErrExpected).Once()

	service := wifi.New(mockWiFi)
	addrs, err := service.GetAddresses()

	require.Error(t, err)
	require.Nil(t, addrs)
	require.Contains(t, err.Error(), "getting interfaces")

	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetNames_Success(t *testing.T) {
	t.Parallel()

	mockWiFi := new(wifi.MockWiFiHandle)

	hwAddr, err := net.ParseMAC("00:11:22:33:44:55")
	require.NoError(t, err)

	interfaces := []*mdlayherwifi.Interface{
		{Name: "wlan0", HardwareAddr: hwAddr},
		{Name: "wlan1", HardwareAddr: nil},
		{Name: "eth0", HardwareAddr: nil},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil).Once()

	service := wifi.New(mockWiFi)
	names, err := service.GetNames()

	require.NoError(t, err)
	require.Len(t, names, 3)
	require.Equal(t, []string{"wlan0", "wlan1", "eth0"}, names)

	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetNames_Empty(t *testing.T) {
	t.Parallel()

	mockWiFi := new(wifi.MockWiFiHandle)
	interfaces := []*mdlayherwifi.Interface{}

	mockWiFi.On("Interfaces").Return(interfaces, nil).Once()

	service := wifi.New(mockWiFi)
	names, err := service.GetNames()

	require.NoError(t, err)
	require.Empty(t, names)

	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetNames_Error(t *testing.T) {
	t.Parallel()

	mockWiFi := new(wifi.MockWiFiHandle)
	mockWiFi.On("Interfaces").Return([]*mdlayherwifi.Interface{}, ErrExpected).Once()

	service := wifi.New(mockWiFi)
	names, err := service.GetNames()

	require.Error(t, err)
	require.Nil(t, names)
	require.Contains(t, err.Error(), "getting interfaces")

	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_BothMethodsSameData(t *testing.T) {
	t.Parallel()

	mockWiFi := new(wifi.MockWiFiHandle)

	hwAddr1, err := net.ParseMAC("00:11:22:33:44:55")
	require.NoError(t, err)
	hwAddr2, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	require.NoError(t, err)
	hwAddr3, err := net.ParseMAC("11:22:33:44:55:66")
	require.NoError(t, err)

	interfaces := []*mdlayherwifi.Interface{
		{Name: "wlan0", HardwareAddr: hwAddr1},
		{Name: "wlan1", HardwareAddr: hwAddr2},
		{Name: "eth0", HardwareAddr: hwAddr3},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil).Twice()

	service := wifi.New(mockWiFi)

	addresses, err := service.GetAddresses()
	require.NoError(t, err, "GetAddresses error")
	require.Len(t, addresses, 3, "Expected 3 addresses")

	names, err := service.GetNames()
	require.NoError(t, err, "GetNames error")
	require.Len(t, names, 3, "Expected 3 names")

	require.Equal(t, "wlan0", names[0], "First interface name mismatch")
	require.Equal(t, "wlan1", names[1], "Second interface name mismatch")
	require.Equal(t, "eth0", names[2], "Third interface name mismatch")

	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_NilHardwareAddr(t *testing.T) {
	t.Parallel()

	mockWiFi := new(wifi.MockWiFiHandle)

	hwAddr, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	require.NoError(t, err)

	interfaces := []*mdlayherwifi.Interface{
		{
			Name:         "wlan0",
			HardwareAddr: nil,
		},
		{Name: "wlan1", HardwareAddr: hwAddr},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil).Once()

	service := wifi.New(mockWiFi)
	addresses, err := service.GetAddresses()
	require.NoError(t, err)
	require.Len(t, addresses, 2, "Expected 2 addresses")
	require.Nil(t, addresses[0], "Expected first address to be nil")
	require.Equal(t, hwAddr.String(), addresses[1].String(), "Second address mismatch")

	mockWiFi.AssertExpectations(t)
}

