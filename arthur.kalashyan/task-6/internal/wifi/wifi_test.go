package wifi_test

import (
	"errors"
	"net"
	"testing"

	mywifi "github.com/Expeline/task-6/internal/wifi"
	wifipkg "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

var errFail = errors.New("fail")

type MockWiFiHandle struct {
	ifaces []*wifipkg.Interface
	err    error
}

func (m *MockWiFiHandle) Interfaces() ([]*wifipkg.Interface, error) {
	return m.ifaces, m.err
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mock    *MockWiFiHandle
		want    []net.HardwareAddr
		wantErr bool
	}{
		{
			name: "single_interface",
			mock: &MockWiFiHandle{
				ifaces: []*wifipkg.Interface{{HardwareAddr: net.HardwareAddr{1, 2, 3}}},
			},
			want: []net.HardwareAddr{{1, 2, 3}},
		},
		{
			name: "multiple_interfaces",
			mock: &MockWiFiHandle{
				ifaces: []*wifipkg.Interface{
					{HardwareAddr: net.HardwareAddr{1}},
					{HardwareAddr: net.HardwareAddr{2}},
				},
			},
			want: []net.HardwareAddr{{1}, {2}},
		},
		{
			name: "empty_interfaces",
			mock: &MockWiFiHandle{ifaces: []*wifipkg.Interface{}},
			want: []net.HardwareAddr{},
		},
		{
			name:    "error",
			mock:    &MockWiFiHandle{err: errFail},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mywifi.New(tt.mock)

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
	t.Parallel()

	tests := []struct {
		name    string
		mock    *MockWiFiHandle
		want    []string
		wantErr bool
	}{
		{
			name: "single_interface",
			mock: &MockWiFiHandle{
				ifaces: []*wifipkg.Interface{{Name: "wlan0"}},
			},
			want: []string{"wlan0"},
		},
		{
			name: "multiple_interfaces",
			mock: &MockWiFiHandle{
				ifaces: []*wifipkg.Interface{
					{Name: "eth0"},
					{Name: "wlan0"},
				},
			},
			want: []string{"eth0", "wlan0"},
		},
		{
			name: "empty_interfaces",
			mock: &MockWiFiHandle{ifaces: []*wifipkg.Interface{}},
			want: []string{},
		},
		{
			name:    "error",
			mock:    &MockWiFiHandle{err: errFail},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mywifi.New(tt.mock)

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
