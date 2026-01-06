package wifi_test

import (
	"errors"
	"net"
	"testing"

	wifipackage "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"

	"github.com/hehemka/task-6/internal/wifi"
)

var errInterfaces = errors.New("interfaces error")

type mockWiFi struct {
	interfaces func() ([]*wifipackage.Interface, error)
}

func (m *mockWiFi) Interfaces() ([]*wifipackage.Interface, error) {
	return m.interfaces()
}

func TestWiFiService_New(t *testing.T) {
	t.Parallel()

	m := &mockWiFi{}
	svc := wifi.New(m)

	require.NotNil(t, svc)
	require.Same(t, m, svc.WiFi)
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	mac1, _ := net.ParseMAC("aa:bb:cc:11:22:33")
	mac2, _ := net.ParseMAC("aa:bb:cc:44:55:66")

	type testCase struct {
		name    string
		mock    *mockWiFi
		want    []net.HardwareAddr
		wantErr bool
	}

	tests := []testCase{
		{
			name: "success",
			mock: &mockWiFi{
				interfaces: func() ([]*wifipackage.Interface, error) {
					return []*wifipackage.Interface{
						{HardwareAddr: mac1},
						{HardwareAddr: mac2},
					}, nil
				},
			},
			want: []net.HardwareAddr{mac1, mac2},
		},
		{
			name: "empty",
			mock: &mockWiFi{
				interfaces: func() ([]*wifipackage.Interface, error) {
					return []*wifipackage.Interface{}, nil
				},
			},
			want: []net.HardwareAddr{},
		},
		{
			name: "error",
			mock: &mockWiFi{
				interfaces: func() ([]*wifipackage.Interface, error) {
					return nil, errInterfaces
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := wifi.New(tt.mock)
			got, err := svc.GetAddresses()

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "getting interfaces")
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name    string
		mock    *mockWiFi
		want    []string
		wantErr bool
	}

	tests := []testCase{
		{
			name: "success",
			mock: &mockWiFi{
				interfaces: func() ([]*wifipackage.Interface, error) {
					return []*wifipackage.Interface{
						{Name: "wlan0"},
						{Name: "wlan1"},
					}, nil
				},
			},
			want: []string{"wlan0", "wlan1"},
		},
		{
			name: "empty",
			mock: &mockWiFi{
				interfaces: func() ([]*wifipackage.Interface, error) {
					return []*wifipackage.Interface{}, nil
				},
			},
			want: []string{},
		},
		{
			name: "error",
			mock: &mockWiFi{
				interfaces: func() ([]*wifipackage.Interface, error) {
					return nil, errInterfaces
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := wifi.New(tt.mock)
			got, err := svc.GetNames()

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "getting interfaces")
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
