package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/danila-clown/task-6/internal/wifi"

	mdlayherWifi "github.com/mdlayher/wifi"
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

	hardwareAddress1, _ := net.ParseMAC("00:11:22:33:44:55")
	hardwareAddress2, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")

	tests := []struct {
		name          string
		interfaces    []*mdlayherWifi.Interface
		mockError     error
		expectedError string
		expectedData  []net.HardwareAddr
	}{
		{
			name: "success - multiple interfaces with MAC addresses",
			interfaces: []*mdlayherWifi.Interface{
				{Name: "wlan0", HardwareAddr: hardwareAddress1},
				{Name: "wlan1", HardwareAddr: hardwareAddress2},
			},
			expectedData: []net.HardwareAddr{hardwareAddress1, hardwareAddress2},
		},
		{
			name:          "error from handle",
			mockError:     errFailedToGetInterfaces,
			expectedError: "getting interfaces: failed to get interfaces",
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

			if tc.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)
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
		name          string
		interfaces    []*mdlayherWifi.Interface
		mockError     error
		expectedError string
		expectedData  []string
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
			name:          "error from handle",
			mockError:     errAccessDenied,
			expectedError: "getting interfaces: access denied",
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

			if tc.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedData, names)
			}
		})
	}
}
