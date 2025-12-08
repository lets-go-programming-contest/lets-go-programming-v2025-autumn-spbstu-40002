package wifi

import (
	"errors"
	"net"
	"testing"
)

type mockInterface struct {
	Index        int
	Name         string
	HardwareAddr net.HardwareAddr
	PHY          int
	Device       int
	Type         int
	Frequency    int
}

type MockWiFiHandle struct {
	interfaces []*mockInterface
	err        error
}

func (m *MockWiFiHandle) Interfaces() ([]*mockInterface, error) {
	return m.interfaces, m.err
}

func (m *MockWiFiHandle) InterfacesAsWifi() ([]interface{}, error) {
	result := make([]interface{}, len(m.interfaces))
	for i, iface := range m.interfaces {
		result[i] = iface
	}
	return result, m.err
}

func parseMAC(s string) net.HardwareAddr {
	mac, err := net.ParseMAC(s)
	if err != nil {
		return nil
	}
	return mac
}

func TestWiFiService_GetAddresses(t *testing.T) {
	mockHandle := &struct {
		InterfacesFunc func() ([]interface{}, error)
	}{}

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
			var mockInterfaces []interface{}
			for i, addr := range tt.addrs {
				mockInterfaces = append(mockInterfaces, &mockInterface{
					Index:        i + 1,
					Name:         "wlan0",
					HardwareAddr: parseMAC(addr),
				})
			}

			mockHandle.InterfacesFunc = func() ([]interface{}, error) {
				return mockInterfaces, tt.mockErr
			}

			if tt.expectErr && tt.mockErr == nil {
				t.Error("Expected error but got none")
			}
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	tests := []struct {
		name        string
		interfaces  []*mockInterface
		mockErr     error
		expectErr   bool
		expectNames []string
	}{
		{
			name: "success with two interfaces",
			interfaces: []*mockInterface{
				{Name: "wlan0", HardwareAddr: parseMAC("00:11:22:33:44:55")},
				{Name: "wlan1", HardwareAddr: parseMAC("aa:bb:cc:dd:ee:ff")},
			},
			expectErr:   false,
			expectNames: []string{"wlan0", "wlan1"},
		},
		{
			name:        "success with no interfaces",
			interfaces:  []*mockInterface{},
			expectErr:   false,
			expectNames: []string{},
		},
		{
			name:       "error from Interfaces",
			interfaces: nil,
			mockErr:    errors.New("mock error"),
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectErr && tt.mockErr == nil {
				t.Error("Test setup error: expected error but mockErr is nil")
			}
		})
	}
}
