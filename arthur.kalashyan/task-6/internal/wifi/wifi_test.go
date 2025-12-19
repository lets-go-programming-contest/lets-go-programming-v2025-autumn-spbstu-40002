package wifi

import (
	"errors"
	"net"
	"testing"

	wifipkg "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

type MockWiFiHandle struct {
	ifaces []*wifipkg.Interface
	err    error
}

func (m *MockWiFiHandle) Interfaces() ([]*wifipkg.Interface, error) {
	return m.ifaces, m.err
}

func TestWiFiService_GetAddresses(t *testing.T) {
	a1 := net.HardwareAddr{1, 2, 3}
	a2 := net.HardwareAddr{4, 5, 6}

	tests := []struct {
		name    string
		mock    *MockWiFiHandle
		want    []net.HardwareAddr
		wantErr bool
	}{
		{
			name: "single_interface",
			mock: &MockWiFiHandle{
				ifaces: []*wifipkg.Interface{
					{HardwareAddr: a1},
				},
			},
			want: []net.HardwareAddr{a1},
		},
		{
			name: "multiple_interfaces",
			mock: &MockWiFiHandle{
				ifaces: []*wifipkg.Interface{
					{HardwareAddr: a1},
					{HardwareAddr: a2},
				},
			},
			want: []net.HardwareAddr{a1, a2},
		},
		{
			name: "empty_interfaces",
			mock: &MockWiFiHandle{
				ifaces: []*wifipkg.Interface{},
			},
			want: []net.HardwareAddr{},
		},
		{
			name: "interfaces_error",
			mock: &MockWiFiHandle{
				err: errors.New("fail"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(tt.mock)

			res, err := svc.GetAddresses()
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, res)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	tests := []struct {
		name    string
		mock    *MockWiFiHandle
		want    []string
		wantErr bool
	}{
		{
			name: "single_interface",
			mock: &MockWiFiHandle{
				ifaces: []*wifipkg.Interface{
					{Name: "wlan0"},
				},
			},
			want: []string{"wlan0"},
		},
		{
			name: "multiple_interfaces_order_preserved",
			mock: &MockWiFiHandle{
				ifaces: []*wifipkg.Interface{
					{Name: "eth0"},
					{Name: "wlan0"},
					{Name: "lo"},
				},
			},
			want: []string{"eth0", "wlan0", "lo"},
		},
		{
			name: "empty_interfaces",
			mock: &MockWiFiHandle{
				ifaces: []*wifipkg.Interface{},
			},
			want: []string{},
		},
		{
			name: "interfaces_error",
			mock: &MockWiFiHandle{
				err: errors.New("fail"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(tt.mock)

			res, err := svc.GetNames()
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, res)
		})
	}
}
