package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mywifi "task-6/internal/wifi"
)

type MockWiFiHandle struct {
	interfaces []*wifi.Interface
	err        error
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	return m.interfaces, m.err
}

func mockInterfaces(addrs []string) []*wifi.Interface {
	var interfaces []*wifi.Interface
	for i, addrStr := range addrs {
		hwAddr, err := net.ParseMAC(addrStr)
		if err != nil {
			continue
		}
		iface := &wifi.Interface{
			Index:        i + 1,
			Name:         fmt.Sprintf("wlan%d", i),
			HardwareAddr: hwAddr,
			PHY:          0,
			Device:       0,
			Type:         wifi.InterfaceTypeStation,
			Frequency:    0,
		}
		interfaces = append(interfaces, iface)
	}
	return interfaces
}

func TestWiFiService_GetAddresses(t *testing.T) {
	tests := []struct {
		name      string
		addrs     []string
		mockErr   error
		expectErr bool
		expectLen int
	}{
		{
			name:      "success with two addresses",
			addrs:     []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
			expectErr: false,
			expectLen: 2,
		},
		{
			name:      "success with no interfaces",
			addrs:     []string{},
			expectErr: false,
			expectLen: 0,
		},
		{
			name:      "error from Interfaces",
			addrs:     nil,
			mockErr:   errors.New("mock error"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockWiFiHandle{
				interfaces: mockInterfaces(tt.addrs),
				err:        tt.mockErr,
			}
			service := mywifi.New(mock)

			addrs, err := service.GetAddresses()
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, addrs, tt.expectLen)

			for i, expectedAddr := range tt.addrs {
				parsed, _ := net.ParseMAC(expectedAddr)
				assert.Equal(t, parsed, addrs[i])
			}
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	tests := []struct {
		name        string
		addrs       []string
		mockErr     error
		expectErr   bool
		expectLen   int
		expectNames []string
	}{
		{
			name:        "success with two interfaces",
			addrs:       []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
			expectErr:   false,
			expectLen:   2,
			expectNames: []string{"wlan0", "wlan1"},
		},
		{
			name:        "success with no interfaces",
			addrs:       []string{},
			expectErr:   false,
			expectLen:   0,
			expectNames: []string{},
		},
		{
			name:        "error from Interfaces",
			addrs:       nil,
			mockErr:     errors.New("mock error"),
			expectErr:   true,
			expectNames: nil,
		},
		{
			name:        "success with one interface",
			addrs:       []string{"11:22:33:44:55:66"},
			expectErr:   false,
			expectLen:   1,
			expectNames: []string{"wlan0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockWiFiHandle{
				interfaces: mockInterfaces(tt.addrs),
				err:        tt.mockErr,
			}
			service := mywifi.New(mock)

			names, err := service.GetNames()
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, names, tt.expectLen)

			if tt.expectNames != nil {
				assert.Equal(t, tt.expectNames, names)
			}
		})
	}
}
