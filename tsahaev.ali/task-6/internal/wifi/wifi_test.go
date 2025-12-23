package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	myWifi "github.com/Tsahaev/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

var errInterfaces = errors.New("interfaces error")

type rowTestSysInfo struct {
	name        string
	addrs       []string
	ifaceNames  []string
	errExpected error
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	service := myWifi.New(mockWifi)
	require.NotNil(t, service)
}

func mockIfaces(addrs []string) []*wifi.Interface {
	if addrs == nil {
		return nil
	}

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

func parseMACs(macStr []string) []net.HardwareAddr {
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

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	testTable := []rowTestSysInfo{
		{
			name:       "Success",
			addrs:      []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
			ifaceNames: []string{"eth1", "eth2"},
		},
		{
			name:        "Error",
			errExpected: errInterfaces,
		},
	}

	for i, row := range testTable {
		t.Run(fmt.Sprintf("%d_%s", i, row.name), func(t *testing.T) {
			t.Parallel()

			mockWifi := NewWiFiHandle(t)
			wifiService := myWifi.New(mockWifi)

			mockWifi.On("Interfaces").Return(mockIfaces(row.addrs), row.errExpected).Once()

			actualAddrs, err := wifiService.GetAddresses()

			if row.errExpected != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, row.errExpected)
				require.Nil(t, actualAddrs)
			} else {
				require.NoError(t, err)
				require.Equal(t, parseMACs(row.addrs), actualAddrs)
			}
		})
	}
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	testTable := []rowTestSysInfo{
		{
			name:       "Success",
			addrs:      []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
			ifaceNames: []string{"eth1", "eth2"},
		},
		{
			name:        "Error",
			errExpected: errInterfaces,
		},
	}

	for i, row := range testTable {
		t.Run(fmt.Sprintf("%d_%s", i, row.name), func(t *testing.T) {
			t.Parallel()

			mockWifi := NewWiFiHandle(t)
			wifiService := myWifi.New(mockWifi)

			mockWifi.On("Interfaces").Return(mockIfaces(row.addrs), row.errExpected).Once()

			actualNames, err := wifiService.GetNames()

			if row.errExpected != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, row.errExpected)
				require.Nil(t, actualNames)
			} else {
				require.NoError(t, err)
				require.Equal(t, row.ifaceNames, actualNames)
			}
		})
	}
}
