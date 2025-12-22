package wifi_test

import (
	"errors"
	"net"
	"testing"

	myWifi "github.com/eg0sha-0/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ErrExpected = errors.New("expected error")

//go:generate mockery --all --testonly --quiet --outpkg wifi_test --output .

func TestWiFiService_GetAddresses_Success(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	interfaces := []*wifi.Interface{
		{Name: "wlan0", HardwareAddr: parseMAC(t, "00:11:22:33:44:55")},
		{Name: "wlan1", HardwareAddr: parseMAC(t, "aa:bb:cc:dd:ee:ff")},
	}
	mockWiFi.On("Interfaces").Return(interfaces, nil)

	addrs, err := service.GetAddresses()
	require.NoError(t, err)
	require.Len(t, addrs, 2)
	assert.Equal(t, "00:11:22:33:44:55", addrs[0].String())
	assert.Equal(t, "aa:bb:cc:dd:ee:ff", addrs[1].String())
	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetAddresses_InterfaceFetchError(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	mockWiFi.On("Interfaces").Return(nil, ErrExpected)

	addrs, err := service.GetAddresses()
	require.ErrorContains(t, err, "getting interfaces:")
	assert.Nil(t, addrs)
	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetAddresses_EmptyList(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, nil)

	addrs, err := service.GetAddresses()
	require.NoError(t, err)
	assert.Empty(t, addrs)
	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetAddresses_NullMACAddress(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	interfaces := []*wifi.Interface{
		{
			Name:         "wlan0",
			HardwareAddr: nil,
		},
	}
	mockWiFi.On("Interfaces").Return(interfaces, nil)

	addrs, err := service.GetAddresses()
	require.NoError(t, err)
	require.Len(t, addrs, 1)
	assert.Nil(t, addrs[0])
	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetNames_Success(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	interfaces := []*wifi.Interface{
		{Name: "wlan0", HardwareAddr: parseMAC(t, "00:11:22:33:44:55")},
		{Name: "eth1", HardwareAddr: parseMAC(t, "11:22:33:44:55:66")},
		{Name: "wifi2", HardwareAddr: parseMAC(t, "22:33:44:55:66:77")},
	}
	mockWiFi.On("Interfaces").Return(interfaces, nil)

	names, err := service.GetNames()
	require.NoError(t, err)
	assert.Equal(t, []string{"wlan0", "eth1", "wifi2"}, names)
	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetNames_InterfaceFetchError(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	mockWiFi.On("Interfaces").Return(nil, ErrExpected)

	names, err := service.GetNames()
	require.ErrorContains(t, err, "getting interfaces:")
	assert.Nil(t, names)
	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetNames_EmptyList(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, nil)

	names, err := service.GetNames()
	require.NoError(t, err)
	assert.Empty(t, names)
	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_GetNames_BlankInterfaceName(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	interfaces := []*wifi.Interface{
		{
			Name:         "",
			HardwareAddr: parseMAC(t, "00:11:22:33:44:55"),
		},
	}
	mockWiFi.On("Interfaces").Return(interfaces, nil)

	names, err := service.GetNames()
	require.NoError(t, err)
	assert.Equal(t, []string{""}, names)
	mockWiFi.AssertExpectations(t)
}

func TestWiFiService_Initialization(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	assert.NotNil(t, service)
	assert.Equal(t, mockWiFi, service.WiFi)
}

func parseMAC(t *testing.T, s string) net.HardwareAddr {
	t.Helper()

	hwAddr, err := net.ParseMAC(s)
	require.NoError(t, err)

	return hwAddr
}
