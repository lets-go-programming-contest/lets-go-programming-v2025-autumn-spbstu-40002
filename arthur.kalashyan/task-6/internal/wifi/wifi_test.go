package wifi_test

import (
	"errors"
	"fmt"
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

	val := args.Get(0)
	if val == nil {
		val = []*wifipkg.Interface{}
	}

	var ifaces []*wifipkg.Interface
	if v, ok := val.([]*wifipkg.Interface); ok {
		ifaces = v
	}

	err := args.Error(1)
	if err != nil {
		err = fmt.Errorf("getting interfaces: %w", err)
	}

	return ifaces, err
}

func TestWiFiService_New(t *testing.T) {
	t.Parallel()

	mockHandle := &MockWiFiHandle{}
	svc := service.New(mockHandle)
	require.NotNil(t, svc)
	require.Same(t, mockHandle, svc.WiFi)
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	mac1, _ := net.ParseMAC("aa:bb:cc:11:22:33")
	mac2, _ := net.ParseMAC("aa:bb:cc:44:55:66")

	tests := []struct {
		name    string
		mock    func() *MockWiFiHandle
		want    []net.HardwareAddr
		wantErr bool
	}{
		{
			name: "multiple",
			mock: func() *MockWiFiHandle {
				m := &MockWiFiHandle{}
				m.On("Interfaces").Return([]*wifipkg.Interface{
					{HardwareAddr: mac1},
					{HardwareAddr: mac2},
				}, nil).Once()

				return m
			},
			want: []net.HardwareAddr{mac1, mac2},
		},
		{
			name: "empty",
			mock: func() *MockWiFiHandle {
				m := &MockWiFiHandle{}
				m.On("Interfaces").Return([]*wifipkg.Interface{}, nil).Once()

				return m
			},
			want: []net.HardwareAddr{},
		},
		{
			name: "error",
			mock: func() *MockWiFiHandle {
				m := &MockWiFiHandle{}
				m.On("Interfaces").Return([]*wifipkg.Interface(nil), errWiFi).Once()

				return m
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := service.New(tt.mock())

			got, err := svc.GetAddresses()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}

			tt.mock().AssertExpectations(t)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mock    func() *MockWiFiHandle
		want    []string
		wantErr bool
	}{
		{
			name: "multiple",
			mock: func() *MockWiFiHandle {
				m := &MockWiFiHandle{}
				m.On("Interfaces").Return([]*wifipkg.Interface{
					{Name: "wlp3s0"}, {Name: "wlan0"},
				}, nil).Once()

				return m
			},
			want: []string{"wlp3s0", "wlan0"},
		},
		{
			name: "empty",
			mock: func() *MockWiFiHandle {
				m := &MockWiFiHandle{}
				m.On("Interfaces").Return([]*wifipkg.Interface{}, nil).Once()

				return m
			},
			want: []string{},
		},
		{
			name: "error",
			mock: func() *MockWiFiHandle {
				m := &MockWiFiHandle{}
				m.On("Interfaces").Return([]*wifipkg.Interface(nil), errWiFi).Once()

				return m
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := service.New(tt.mock())

			got, err := svc.GetNames()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}

			tt.mock().AssertExpectations(t)
		})
	}
}
