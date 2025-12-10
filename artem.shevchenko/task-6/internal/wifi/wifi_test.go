package wifi_test

import (
	"errors"
	"net"
	"testing"

	mdlayherWifi "github.com/mdlayher/wifi"
	"github.com/slendycs/go-lab-6/internal/wifi"
	"github.com/stretchr/testify/require"
)

var (
	errFailedToGetInterfaces = errors.New("failed to get interfaces")
	errAccessDenied          = errors.New("access denied")
)

type mockWiFiHandle struct {
	interfaces []*mdlayherWifi.Interface
	err        error
}

func (m *mockWiFiHandle) Interfaces() ([]*mdlayherWifi.Interface, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.interfaces, nil
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	hwAddr1, _ := net.ParseMAC("00:11:22:33:44:55")
	hwAddr2, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")

	tests := []struct {
		name         string
		interfaces   []*mdlayherWifi.Interface
		mockError    error
		expectedErr  string
		expectedData []net.HardwareAddr
	}{
		{
			name: "success - multiple interfaces with MAC addresses",
			interfaces: []*mdlayherWifi.Interface{
				{Name: "wlan0", HardwareAddr: hwAddr1},
				{Name: "wlan1", HardwareAddr: hwAddr2},
			},
			expectedData: []net.HardwareAddr{hwAddr1, hwAddr2},
		},
		{
			name:        "error from handle",
			mockError:   errFailedToGetInterfaces,
			expectedErr: "getting interfaces: failed to get interfaces",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandle := &mockWiFiHandle{
				interfaces: tc.interfaces,
				err:        tc.mockError,
			}

			service := wifi.New(mockHandle)
			addrs, err := service.GetAddresses()

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
				require.Nil(t, addrs)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedData, addrs)
			}
		})
	}
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	hwAddr, _ := net.ParseMAC("00:11:22:33:44:55")

	tests := []struct {
		name         string
		interfaces   []*mdlayherWifi.Interface
		mockError    error
		expectedErr  string
		expectedData []string
	}{
		{
			name: "success - multiple interfaces",
			interfaces: []*mdlayherWifi.Interface{
				{Name: "wlan0", HardwareAddr: hwAddr},
				{Name: "wlan1", HardwareAddr: hwAddr},
			},
			expectedData: []string{"wlan0", "wlan1"},
		},
		{
			name:        "error from handle",
			mockError:   errAccessDenied,
			expectedErr: "getting interfaces: access denied",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandle := &mockWiFiHandle{
				interfaces: tc.interfaces,
				err:        tc.mockError,
			}

			service := wifi.New(mockHandle)
			names, err := service.GetNames()

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedData, names)
			}
		})
	}
}
