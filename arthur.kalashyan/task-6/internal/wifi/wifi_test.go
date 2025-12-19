package wifi_test

import (
	"errors"
	"net"
	"testing"

	service "github.com/Expeline/task-6/internal/wifi"
	wifipkg "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errWiFi = errors.New("failed to get interfaces")

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifipkg.Interface, error) {
	args := m.Called()
	ifaces := args.Get(0)
	if ifaces != nil {
		if val, ok := ifaces.([]*wifipkg.Interface); ok {
			ifaces = val
		}
	}
	return ifaces.([]*wifipkg.Interface), args.Error(1)
}

func TestWiFiService_New(t *testing.T) {
	t.Parallel()
	mockHandle := &MockWiFiHandle{}
	svc := service.New(mockHandle)
	require.NotNil(t, svc)
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
			name: "multiple",
			mock: func() *MockWiFiHandle {
				m := &MockWiFiHandle{}
				mac1, _ := net.ParseMAC("aa:bb:cc:11:22:33")
				mac2, _ := net.ParseMAC("aa:bb:cc:44:55:66")
				m.On("Interfaces").Return([]*wifipkg.Interface{
					{HardwareAddr: mac1},
					{HardwareAddr: mac2},
				}, nil).Once()
				return m
			}(),
			want: []net.HardwareAddr{
				{0xaa, 0xbb, 0xcc, 0x11, 0x22, 0x33},
				{0xaa, 0xbb, 0xcc, 0x44, 0x55, 0x66},
			},
		},
		{
			name: "empty",
			mock: func() *MockWiFiHandle {
				m := &MockWiFiHandle{}
				m.On("Interfaces").Return([]*wifipkg.Interface{}, nil).Once()
				return m
			}(),
			want: []net.HardwareAddr{},
		},
		{
			name: "error",
			mock: func() *MockWiFiHandle {
				m := &MockWiFiHandle{}
				m.On("Interfaces").Return([]*wifipkg.Interface(nil), errWiFi).Once()
				return m
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := service.New(tt.mock)
			got, err := svc.GetAddresses()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
			tt.mock.AssertExpectations(t)
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
			name: "multiple",
			mock: func() *MockWiFiHandle {
				m := &MockWiFiHandle{}
				m.On("Interfaces").Return([]*wifipkg.Interface{
					{Name: "wlp3s0"},
					{Name: "wlan0"},
				}, nil).Once()
				return m
			}(),
			want: []string{"wlp3s0", "wlan0"},
		},
		{
			name: "empty",
			mock: func() *MockWiFiHandle {
				m := &MockWiFiHandle{}
				m.On("Interfaces").Return([]*wifipkg.Interface{}, nil).Once()
				return m
			}(),
			want: []string{},
		},
		{
			name: "error",
			mock: func() *MockWiFiHandle {
				m := &MockWiFiHandle{}
				m.On("Interfaces").Return([]*wifipkg.Interface(nil), errWiFi).Once()
				return m
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := service.New(tt.mock)
			got, err := svc.GetNames()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
			tt.mock.AssertExpectations(t)
		})
	}
}
