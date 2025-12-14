package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	myWifi "github.com/HuaChenju/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

type rowTestSysInfo struct {
	addrs       []string
	names       []string
	errExpected error
}

var testTable = []rowTestSysInfo{ //nolint:gochecknoglobals
	{
		addrs: []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
		names: []string{"eth1", "eth2"},
	},
	{
		addrs:       nil,
		errExpected: ErrExpected,
	},
	{
		addrs: []string{},
		names: []string{},
	},
}

var ErrExpected = errors.New("expected error")

//go:generate mockery --all --testonly --quiet --outpkg wifi_test --output .

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	wifiService := myWifi.WiFiService{WiFi: mockWifi}

	for i, row := range testTable {
		mockWifi.On("Interfaces").Unset()
		mockWifi.On("Interfaces").Return(helperMockIfaces(t, row.addrs), row.errExpected)

		actualAddrs, err := wifiService.GetAddresses()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected)

			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", i)
		require.Equal(t, helperParseMACs(t, row.addrs), actualAddrs,
			"row: %d, expected addrs: %s, actual addrs: %s", i,
			helperParseMACs(t, row.addrs), actualAddrs)
	}

	mockWifi.AssertExpectations(t)
}

func helperMockIfaces(t *testing.T, addrs []string) []*wifi.Interface {
	t.Helper()

	interfaces := make([]*wifi.Interface, 0, len(addrs))

	for i, addrStr := range addrs {
		hwAddr := parseMAC(addrStr)
		if hwAddr == nil {
			continue
		}

		iface := &wifi.Interface{
			Index:        i + 1,
			Name:         fmt.Sprintf("eth%d", i+1),
			HardwareAddr: hwAddr,
			PHY:          1,
			Device:       1,
			Type:         wifi.InterfaceTypeAPVLAN,
			Frequency:    0,
		}

		interfaces = append(interfaces, iface)
	}

	return interfaces
}

func helperParseMACs(t *testing.T, macStr []string) []net.HardwareAddr {
	t.Helper()

	addrs := make([]net.HardwareAddr, 0, len(macStr))

	for _, addr := range macStr {
		addrs = append(addrs, parseMAC(addr))
	}

	return addrs
}

func parseMAC(macStr string) net.HardwareAddr {
	hwAddr, err := net.ParseMAC(macStr)
	if err != nil {
		return nil
	}

	return hwAddr
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	wifiService := myWifi.WiFiService{WiFi: mockWifi}

	for i, row := range testTable {
		mockWifi.On("Interfaces").Unset()
		mockWifi.On("Interfaces").Return(helperMockIfaces(t, row.addrs), row.errExpected)

		actualNames, err := wifiService.GetNames()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected)

			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", i)
		require.Equal(t, row.names, actualNames, "row: %d, expected names: %v, actual names: %v",
			i, row.names, actualNames,
		)
	}

	mockWifi.AssertExpectations(t)
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	wifiService := myWifi.New(mockWifi)
	require.Equal(t, mockWifi, wifiService.WiFi)
}
