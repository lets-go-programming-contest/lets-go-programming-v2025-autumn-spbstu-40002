package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockWiFiHandle struct {
	mock.Mock
}

func (m *mockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*wifi.Interface), args.Error(1)
}

type rowTestSysInfo1 struct {
	addrs       []string
	errExpected error
}

var testTable2 = []rowTestSysInfo1{
	{
		addrs: []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
	},
	{
		errExpected: errors.New("ExpectedError"),
	},
	{
		addrs: []string{}, // тест пустого списка
	},
}

type rowTestSysInfo2 struct {
	names       []string
	errExpected error
}

var testTable3 = []rowTestSysInfo2{
	{
		names: []string{"wlan0", "eth0", "wlp3s0"},
	},
	{
		errExpected: errors.New("ExpectedError"),
	},
	{
		names: []string{}, // тест пустого списка
	},
}

func TestGetAddresses(t *testing.T) {
	for i, row := range testTable2 {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			mockWifi := new(mockWiFiHandle)

			if row.errExpected != nil {
				mockWifi.On("Interfaces").Return(nil, row.errExpected)
			} else {
				interfaces := mockIfacesWithAdds(row.addrs)
				mockWifi.On("Interfaces").Return(interfaces, nil)
			}

			wifiService := New(mockWifi)
			actualAddrs, err := wifiService.GetAddresses()

			if row.errExpected != nil {
				require.Error(t, err, "row: %d, expected error but got none", i)
				require.ErrorContains(t, err, row.errExpected.Error(), "row: %d", i)
				require.Nil(t, actualAddrs, "row: %d, addrs must be nil", i)
			} else {
				require.NoError(t, err, "row: %d, error must be nil", i)
				expectedAddrs := parseMACs(row.addrs)

				if expectedAddrs == nil && actualAddrs != nil {
					require.Len(t, actualAddrs, 0,
						"row: %d, expected nil or empty slice, got len=%d", i, len(actualAddrs))
				} else {
					require.Equal(t, expectedAddrs, actualAddrs,
						"row: %d, expected addrs: %v actual addrs: %v", i, expectedAddrs, actualAddrs)
				}
			}

			mockWifi.AssertExpectations(t)
		})
	}
}

func TestGetNames(t *testing.T) {
	for i, row := range testTable3 {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			mockWifi := new(mockWiFiHandle)

			if row.errExpected != nil {
				mockWifi.On("Interfaces").Return(nil, row.errExpected)
			} else {
				interfaces := mockIfacesWithNames(row.names)
				mockWifi.On("Interfaces").Return(interfaces, nil)
			}

			wifiService := New(mockWifi)
			actualNames, err := wifiService.GetNames()

			if row.errExpected != nil {
				require.Error(t, err, "row: %d, expected error but got none", i)
				require.ErrorContains(t, err, row.errExpected.Error(), "row: %d", i)
				require.Nil(t, actualNames, "row: %d, names must be nil", i)
			} else {
				require.NoError(t, err, "row: %d, error must be nil", i)
				require.Equal(t, row.names, actualNames,
					"row: %d, expected names: %v, actual names: %v", i, row.names, actualNames)
			}

			mockWifi.AssertExpectations(t)
		})
	}
}

func mockIfacesWithAdds(addrs []string) []*wifi.Interface {
	var interfaces []*wifi.Interface
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

func mockIfacesWithNames(names []string) []*wifi.Interface {
	var interfaces []*wifi.Interface
	for i, name := range names {
		iface := &wifi.Interface{
			Index:        i + 1,
			Name:         name,
			HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
			PHY:          1,
			Device:       1,
			Type:         wifi.InterfaceTypeAPVLAN,
			Frequency:    0,
		}
		interfaces = append(interfaces, iface)
	}
	return interfaces
}

func parseMACs(macStrs []string) []net.HardwareAddr {
	var addrs []net.HardwareAddr
	for _, addrStr := range macStrs {
		hwAddr := parseMAC(addrStr)
		if hwAddr != nil {
			addrs = append(addrs, hwAddr)
		}
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

