package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/XShaygaND/task-6/internal/wifi"
	mdlayherwifi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errIFace = errors.New("interface retrieval error")

func TestCreateManager(t *testing.T) {
	t.Parallel()

	mockHandler := &wifi.MockWiFiHandle{}
	manager := wifi.CreateManager(mockHandler)

	assert.Equal(t, mockHandler, manager.Handler)
}

func TestWiFiManager_GetMACAddresses(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		setup    func(mock *wifi.MockWiFiHandle)
		expected []net.HardwareAddr
		errMsg   string
	}{
		{
			name: "multiple interfaces",
			setup: func(mock *wifi.MockWiFiHandle) {
				ifaces := []*mdlayherwifi.Interface{
					{HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}},
					{HardwareAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
				}
				mock.On("Interfaces").Return(ifaces, nil)
			},
			expected: []net.HardwareAddr{
				{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			},
		},
		{
			name: "no interfaces",
			setup: func(mock *wifi.MockWiFiHandle) {
				mock.On("Interfaces").Return([]*mdlayherwifi.Interface{}, nil)
			},
			expected: []net.HardwareAddr{},
		},
		{
			name: "interface error",
			setup: func(mock *wifi.MockWiFiHandle) {
				mock.On("Interfaces").Return([]*mdlayherwifi.Interface(nil), errIFace)
			},
			errMsg: "failed to retrieve interfaces: interface retrieval error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandler := &wifi.MockWiFiHandle{}
			tc.setup(mockHandler)

			manager := wifi.CreateManager(mockHandler)
			result, err := manager.GetMACAddresses()

			if tc.errMsg != "" {
				require.Error(t, err)
				assert.Equal(t, tc.errMsg, err.Error())
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}

			mockHandler.AssertExpectations(t)
		})
	}
}

func TestWiFiManager_GetInterfaceNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		setup    func(mock *wifi.MockWiFiHandle)
		expected []string
		errMsg   string
	}{
		{
			name: "multiple interface names",
			setup: func(mock *wifi.MockWiFiHandle) {
				ifaces := []*mdlayherwifi.Interface{
					{Name: "eth0"},
					{Name: "wlan0"},
				}
				mock.On("Interfaces").Return(ifaces, nil)
			},
			expected: []string{"eth0", "wlan0"},
		},
		{
			name: "no interface names",
			setup: func(mock *wifi.MockWiFiHandle) {
				mock.On("Interfaces").Return([]*mdlayherwifi.Interface{}, nil)
			},
			expected: []string{},
		},
		{
			name: "interface name error",
			setup: func(mock *wifi.MockWiFiHandle) {
				mock.On("Interfaces").Return([]*mdlayherwifi.Interface(nil), errIFace)
			},
			errMsg: "failed to retrieve interfaces: interface retrieval error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandler := &wifi.MockWiFiHandle{}
			tc.setup(mockHandler)

			manager := wifi.CreateManager(mockHandler)
			result, err := manager.GetInterfaceNames()

			if tc.errMsg != "" {
				require.Error(t, err)
				assert.Equal(t, tc.errMsg, err.Error())
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}

			mockHandler.AssertExpectations(t)
		})
	}
}
