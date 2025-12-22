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

var errInterfaces = errors.New("interface retrieval error")

func TestNew(t *testing.T) {
	t.Parallel()

	mockHandle := &wifi.MockWiFiHandle{}
	s := wifi.New(mockHandle)

	assert.Equal(t, mockHandle, s.WiFi)
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		setup    func(mock *wifi.MockWiFiHandle)
		expected []net.HardwareAddr
		errMsg   string
	}{
		{
			name: "success multiple interfaces",
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
			name: "success empty",
			setup: func(mock *wifi.MockWiFiHandle) {
				mock.On("Interfaces").Return([]*mdlayherwifi.Interface{}, nil)
			},
			expected: []net.HardwareAddr{},
		},
		{
			name: "error",
			setup: func(mock *wifi.MockWiFiHandle) {
				mock.On("Interfaces").Return([]*mdlayherwifi.Interface(nil), errInterfaces)
			},
			errMsg: "failed to retrieve interfaces: interface retrieval error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandle := &wifi.MockWiFiHandle{}
			tc.setup(mockHandle)

			s := wifi.New(mockHandle)
			got, err := s.GetAddresses()

			if tc.errMsg != "" {
				require.Error(t, err)

				assert.Equal(t, tc.errMsg, err.Error())
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)

				assert.Equal(t, tc.expected, got)
			}

			mockHandle.AssertExpectations(t)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		setup    func(mock *wifi.MockWiFiHandle)
		expected []string
		errMsg   string
	}{
		{
			name: "success multiple interfaces",
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
			name: "success empty",
			setup: func(mock *wifi.MockWiFiHandle) {
				mock.On("Interfaces").Return([]*mdlayherwifi.Interface{}, nil)
			},
			expected: []string{},
		},
		{
			name: "error",
			setup: func(mock *wifi.MockWiFiHandle) {
				mock.On("Interfaces").Return([]*mdlayherwifi.Interface(nil), errInterfaces)
			},
			errMsg: "failed to retrieve interfaces: interface retrieval error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandle := &wifi.MockWiFiHandle{}
			tc.setup(mockHandle)

			s := wifi.New(mockHandle)
			got, err := s.GetNames()

			if tc.errMsg != "" {
				require.Error(t, err)

				assert.Equal(t, tc.errMsg, err.Error())
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)

				assert.Equal(t, tc.expected, got)
			}

			mockHandle.AssertExpectations(t)
		})
	}
}
