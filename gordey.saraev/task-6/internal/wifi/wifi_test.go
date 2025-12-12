package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/F0LY/task-6/internal/wifi"
	mdlayherWifi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*mdlayherWifi.Interface, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*mdlayherWifi.Interface), args.Error(1)
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		mockHandle := &MockWiFiHandle{}

		service := wifi.New(mockHandle)

		ifaces := []*mdlayherWifi.Interface{
			{
				HardwareAddr: net.HardwareAddr{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
			},
			{
				HardwareAddr: net.HardwareAddr{0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b},
			},
		}

		mockHandle.On("Interfaces").Return(ifaces, nil)

		addrs, err := service.GetAddresses()

		require.NoError(t, err)
		require.Equal(t, []net.HardwareAddr{ifaces[0].HardwareAddr, ifaces[1].HardwareAddr}, addrs)

		mockHandle.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockHandle := &MockWiFiHandle{}

		service := wifi.New(mockHandle)

		mockHandle.On("Interfaces").Return(nil, errors.New("interfaces error"))

		_, err := service.GetAddresses()

		require.Error(t, err)

		mockHandle.AssertExpectations(t)
	})

	t.Run("empty", func(t *testing.T) {
		mockHandle := &MockWiFiHandle{}

		service := wifi.New(mockHandle)

		mockHandle.On("Interfaces").Return([]*mdlayherWifi.Interface{}, nil)

		addrs, err := service.GetAddresses()

		require.NoError(t, err)
		require.Empty(t, addrs)

		mockHandle.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		mockHandle := &MockWiFiHandle{}

		service := wifi.New(mockHandle)

		ifaces := []*mdlayherWifi.Interface{
			{
				Name: "wlan0",
			},
			{
				Name: "wlan1",
			},
		}

		mockHandle.On("Interfaces").Return(ifaces, nil)

		names, err := service.GetNames()

		require.NoError(t, err)
		require.Equal(t, []string{"wlan0", "wlan1"}, names)

		mockHandle.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockHandle := &MockWiFiHandle{}

		service := wifi.New(mockHandle)

		mockHandle.On("Interfaces").Return(nil, errors.New("interfaces error"))

		_, err := service.GetNames()

		require.Error(t, err)

		mockHandle.AssertExpectations(t)
	})

	t.Run("empty", func(t *testing.T) {
		mockHandle := &MockWiFiHandle{}

		service := wifi.New(mockHandle)

		mockHandle.On("Interfaces").Return([]*mdlayherWifi.Interface{}, nil)

		names, err := service.GetNames()

		require.NoError(t, err)
		require.Empty(t, names)

		mockHandle.AssertExpectations(t)
	})
}
