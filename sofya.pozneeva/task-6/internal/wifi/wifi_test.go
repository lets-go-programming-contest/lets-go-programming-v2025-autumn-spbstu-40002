package wifi_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	myWifi "github.com/sonychello/task-6/internal/wifi"
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

func TestGetAddresses(t *testing.T) {
	tests := []struct {
		name        string
		addrs       []string
		errExpected error
	}{
		{
			name:  "multiple addresses",
			addrs: []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
		},
		{
			name:        "error from interface",
			errExpected: fmt.Errorf("ExpectedError"),
		},
		{
			name:  "no addresses",
			addrs: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockWifi := new(mockWiFiHandle)
			if tc.errExpected != nil {
				mockWifi.On("Interfaces").Return(nil, tc.errExpected)
			} else {
				interfaces := mockIfacesWithAdds(tc.addrs)
				mockWifi.On("Interfaces").Return(interfaces, nil)
			}

			wifiService := myWifi.New(mockWifi)
			actualAddrs, err := wifiService.GetAddresses()

			if tc.errExpected != nil {
				require.Error(t, err)
				require.Nil(t, actualAddrs)
			} else {
				require.NoError(t, err)
				expectedAddrs := parseMACs(tc.addrs)
				require.Equal(t, expectedAddrs, actualAddrs)
			}
			mockWifi.AssertExpectations(t)
		})
	}
}

func TestGetNames(t *testing.T) {
	tests := []struct {
		name        string
		names       []string
		errExpected error
	}{
		{
			name:  "multiple names",
			names: []string{"wlan0", "eth0", "wlp3s0"},
		},
		{
			name:        "error from interface",
			errExpected: fmt.Errorf("ExpectedError"),
		},
		{
			name:  "no names",
			names: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockWifi := new(mockWiFiHandle)
			if tc.errExpected != nil {
				mockWifi.On("Interfaces").Return(nil, tc.errExpected)
			} else {
				interfaces := mockIfacesWithNames(tc.names)
				mockWifi.On("Interfaces").Return(interfaces, nil)
			}

			wifiService := myWifi.New(mockWifi)
			actualNames, err := wifiService.GetNames()

			if tc.errExpected != nil {
				require.Error(t, err)
				require.Nil(t, actualNames)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.names, actualNames)
			}
			mockWifi.AssertExpectations(t)
		})
	}
}

func mockIfacesWithAdds(addrs []string) []*wifi.Interface {
	var interfaces []*wifi.Interface
	for i, addrStr := range addrs {
		hwAddr, _ := net.ParseMAC(addrStr)
		if hwAddr == nil {
			continue
		}
		iface := &wifi.Interface{
			Index:        i + 1,
			Name:         fmt.Sprintf("eth%d", i+1),
			HardwareAddr: hwAddr,
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
		}
		interfaces = append(interfaces, iface)
	}
	return interfaces
}

func parseMACs(macStrs []string) []net.HardwareAddr {
	var addrs []net.HardwareAddr
	for _, addrStr := range macStrs {
		hwAddr, _ := net.ParseMAC(addrStr)
		if hwAddr != nil {
			addrs = append(addrs, hwAddr)
		}
	}
	return addrs
}

