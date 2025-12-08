package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockWiFiHandle struct {
	InterfacesFunc func() ([]*wifi.Interface, error)
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	if m.InterfacesFunc != nil {
		return m.InterfacesFunc()
	}
	return nil, nil
}

func parseMAC(s string) net.HardwareAddr {
	mac, err := net.ParseMAC(s)
	if err != nil {
		return nil
	}
	return mac
}

func TestWiFiService_GetAddresses(t *testing.T) {
	mac1 := parseMAC("00:11:22:33:44:55")
	mac2 := parseMAC("aa:bb:cc:dd:ee:ff")

	tests := []struct {
		name         string
		interfaces   []*wifi.Interface
		mockErr      error
		expectErr    bool
		expectResult []net.HardwareAddr
	}{
		{
			name: "success with multiple interfaces",
			interfaces: []*wifi.Interface{
				{Name: "wlan0", HardwareAddr: mac1},
				{Name: "wlan1", HardwareAddr: mac2},
			},
			expectResult: []net.HardwareAddr{mac1, mac2},
		},
		{
			name:         "error from Interfaces",
			mockErr:      errors.New("mock error"),
			expectErr:    true,
			expectResult: nil,
		},
		{
			name:         "empty interfaces",
			interfaces:   []*wifi.Interface{},
			expectResult: []net.HardwareAddr{},
		},
		{
			name: "interface with nil hardware address",
			interfaces: []*wifi.Interface{
				{Name: "wlan0", HardwareAddr: nil},
			},
			expectResult: []net.HardwareAddr{nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockWiFiHandle{
				InterfacesFunc: func() ([]*wifi.Interface, error) {
					return tt.interfaces, tt.mockErr
				},
			}

			service := New(mock)
			addrs, err := service.GetAddresses()

			if tt.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "getting interfaces")
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectResult, addrs)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	mac1 := parseMAC("00:11:22:33:44:55")
	mac2 := parseMAC("aa:bb:cc:dd:ee:ff")

	tests := []struct {
		name         string
		interfaces   []*wifi.Interface
		mockErr      error
		expectErr    bool
		expectResult []string
	}{
		{
			name: "success with multiple interfaces",
			interfaces: []*wifi.Interface{
				{Name: "wlan0", HardwareAddr: mac1},
				{Name: "wlan1", HardwareAddr: mac2},
			},
			expectResult: []string{"wlan0", "wlan1"},
		},
		{
			name:         "error from Interfaces",
			mockErr:      errors.New("mock error"),
			expectErr:    true,
			expectResult: nil,
		},
		{
			name:         "empty interfaces",
			interfaces:   []*wifi.Interface{},
			expectResult: []string{},
		},
		{
			name: "interface with empty name",
			interfaces: []*wifi.Interface{
				{Name: "", HardwareAddr: mac1},
			},
			expectResult: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockWiFiHandle{
				InterfacesFunc: func() ([]*wifi.Interface, error) {
					return tt.interfaces, tt.mockErr
				},
			}

			service := New(mock)
			names, err := service.GetNames()

			if tt.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "getting interfaces")
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectResult, names)
		})
	}
}

func TestWiFiService_New(t *testing.T) {
	mockWiFi := &MockWiFiHandle{}
	service := New(mockWiFi)
	assert.Equal(t, mockWiFi, service.WiFi)
}
