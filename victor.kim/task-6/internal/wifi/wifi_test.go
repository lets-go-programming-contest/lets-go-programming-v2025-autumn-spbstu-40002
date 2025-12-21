package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	service "github.com/lets-go-programming-v2025-autumn-spbstu-40002/victor.kim/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (_m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	ret := _m.Called()

	var r0 []*wifi.Interface

	if val := ret.Get(0); val != nil {
		if v, ok := val.([]*wifi.Interface); ok {
			r0 = v
		}
	}

	var r1 error
	if err := ret.Error(1); err != nil {
		r1 = fmt.Errorf("mock Interfaces error: %w", err)
	}

	return r0, r1
}

var errWiFi = errors.New("failed to get interfaces")

func TestWiFiService_New(t *testing.T) {
	t.Parallel()

	mockHandle := &MockWiFiHandle{}
	svc := service.New(mockHandle)

	assert.NotNil(t, svc)
	assert.Same(t, mockHandle, svc.WiFi)
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success - multiple interfaces", func(t *testing.T) {
		t.Parallel()

		mockHandle := &MockWiFiHandle{}
		svc := service.New(mockHandle)

		mac1, _ := net.ParseMAC("aa:bb:cc:00:00:01")
		mac2, _ := net.ParseMAC("aa:bb:cc:00:00:02")

		ifaces := []*wifi.Interface{
			{HardwareAddr: mac1},
			{HardwareAddr: mac2},
		}

		mockHandle.On("Interfaces").Return(ifaces, nil).Once()

		addrs, err := svc.GetAddresses()

		require.NoError(t, err)
		assert.Equal(t, []net.HardwareAddr{mac1, mac2}, addrs)
		mockHandle.AssertExpectations(t)
	})

	t.Run("success - empty", func(t *testing.T) {
		t.Parallel()

		mockHandle := &MockWiFiHandle{}
		svc := service.New(mockHandle)

		mockHandle.On("Interfaces").Return([]*wifi.Interface{}, nil).Once()

		addrs, err := svc.GetAddresses()

		require.NoError(t, err)
		assert.Empty(t, addrs)
		mockHandle.AssertExpectations(t)
	})

	t.Run("error from Interfaces", func(t *testing.T) {
		t.Parallel()

		mockHandle := &MockWiFiHandle{}
		svc := service.New(mockHandle)

		mockHandle.On("Interfaces").Return([]*wifi.Interface(nil), errWiFi).Once()

		addrs, err := svc.GetAddresses()

		require.Error(t, err)
		assert.Nil(t, addrs)
		assert.Contains(t, err.Error(), "getting interfaces")
		mockHandle.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success - multiple names", func(t *testing.T) {
		t.Parallel()

		mockHandle := &MockWiFiHandle{}
		svc := service.New(mockHandle)

		ifaces := []*wifi.Interface{
			{Name: "wlp3s0"},
			{Name: "wlan0"},
		}

		mockHandle.On("Interfaces").Return(ifaces, nil).Once()

		names, err := svc.GetNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"wlp3s0", "wlan0"}, names)
		mockHandle.AssertExpectations(t)
	})

	t.Run("success - empty", func(t *testing.T) {
		t.Parallel()

		mockHandle := &MockWiFiHandle{}
		svc := service.New(mockHandle)

		mockHandle.On("Interfaces").Return([]*wifi.Interface{}, nil).Once()

		names, err := svc.GetNames()

		require.NoError(t, err)
		assert.Empty(t, names)
		mockHandle.AssertExpectations(t)
	})

	t.Run("error from Interfaces", func(t *testing.T) {
		t.Parallel()

		mockHandle := &MockWiFiHandle{}
		svc := service.New(mockHandle)

		mockHandle.On("Interfaces").Return([]*wifi.Interface(nil), errWiFi).Once()

		names, err := svc.GetNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "getting interfaces")
		mockHandle.AssertExpectations(t)
	})
}
