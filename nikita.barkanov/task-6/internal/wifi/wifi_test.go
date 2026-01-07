package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	customWifi "github.com/ControlShiftEscape/task-6/internal/wifi"
)

var (
	errPermissionDenied = errors.New("permission denied")
	errIOTimeout        = errors.New("i/o timeout")
)

type mockWiFiHandle struct {
	interfacesFunc func() ([]*wifi.Interface, error)
}

func (m *mockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	if m.interfacesFunc != nil {
		return m.interfacesFunc()
	}

	return nil, nil
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		mockSetup     func() *mockWiFiHandle
		wantAddrs     []net.HardwareAddr
		wantErr       bool
		wantErrString string
	}{
		{
			name: "success_multiple_macs",
			mockSetup: func() *mockWiFiHandle {
				return &mockWiFiHandle{
					interfacesFunc: func() ([]*wifi.Interface, error) {
						return []*wifi.Interface{
							{Name: "wlan0", HardwareAddr: mac("00:11:22:33:44:55")},
							{Name: "wlan1", HardwareAddr: mac("aa:bb:cc:dd:ee:ff")},
							{Name: "wlp3s0", HardwareAddr: mac("01:02:03:04:05:06")},
						}, nil
					},
				}
			},
			wantAddrs: []net.HardwareAddr{
				mac("00:11:22:33:44:55"),
				mac("aa:bb:cc:dd:ee:ff"),
				mac("01:02:03:04:05:06"),
			},
			wantErr: false,
		},
		{
			name: "empty_interfaces_list",
			mockSetup: func() *mockWiFiHandle {
				return &mockWiFiHandle{
					interfacesFunc: func() ([]*wifi.Interface, error) {
						return []*wifi.Interface{}, nil
					},
				}
			},
			wantAddrs: []net.HardwareAddr{},
			wantErr:   false,
		},
		{
			name: "error_from_wifi_handle",
			mockSetup: func() *mockWiFiHandle {
				return &mockWiFiHandle{
					interfacesFunc: func() ([]*wifi.Interface, error) {
						return nil, errPermissionDenied
					},
				}
			},
			wantAddrs:     nil,
			wantErr:       true,
			wantErrString: "getting interfaces: permission denied",
		},
		{
			name: "interface_with_nil_hardware_addr",
			mockSetup: func() *mockWiFiHandle {
				return &mockWiFiHandle{
					interfacesFunc: func() ([]*wifi.Interface, error) {
						return []*wifi.Interface{
							{Name: "wlan0", HardwareAddr: nil},
							{Name: "wlan1", HardwareAddr: mac("aa:bb:cc:dd:ee:ff")},
						}, nil
					},
				}
			},
			wantAddrs: []net.HardwareAddr{
				nil,
				mac("aa:bb:cc:dd:ee:ff"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mock := tt.mockSetup()
			service := customWifi.New(mock)

			got, err := service.GetAddresses()

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrString)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantAddrs, got)
			}
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		mockSetup func() *mockWiFiHandle
		wantNames []string
		wantErr   bool
	}{
		{
			name: "multiple_interfaces",
			mockSetup: func() *mockWiFiHandle {
				return &mockWiFiHandle{
					interfacesFunc: func() ([]*wifi.Interface, error) {
						return []*wifi.Interface{
							{Name: "wlan0"},
							{Name: "wlan1"},
							{Name: "ap0"},
						}, nil
					},
				}
			},
			wantNames: []string{"wlan0", "wlan1", "ap0"},
			wantErr:   false,
		},
		{
			name: "empty_list",
			mockSetup: func() *mockWiFiHandle {
				return &mockWiFiHandle{
					interfacesFunc: func() ([]*wifi.Interface, error) {
						return nil, nil
					},
				}
			},
			wantNames: []string{},
			wantErr:   false,
		},
		{
			name: "error_on_get_interfaces",
			mockSetup: func() *mockWiFiHandle {
				return &mockWiFiHandle{
					interfacesFunc: func() ([]*wifi.Interface, error) {
						return nil, errIOTimeout
					},
				}
			},
			wantNames: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mock := tt.mockSetup()
			service := customWifi.New(mock)

			got, err := service.GetNames()

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantNames, got)
			}
		})
	}
}

func mac(s string) net.HardwareAddr {
	addr, err := net.ParseMAC(s)
	if err != nil {
		panic(err)
	}

	return addr
}
