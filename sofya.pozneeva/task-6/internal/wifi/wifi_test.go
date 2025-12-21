package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"

	wifisvc "task-6/internal/wifi"
)

var errBoom = errors.New("boom")

type handleMock struct {
	ifaces []*wifi.Interface
	err    error
	calls  int
}

func (h *handleMock) Interfaces() ([]*wifi.Interface, error) {
	h.calls++

	return h.ifaces, h.err
}

func TestWiFiService_GetAddresses_OK(t *testing.T) {
	t.Parallel()

	h := &handleMock{
		ifaces: []*wifi.Interface{
			{
				Name:         "wlan0",
				HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
			},
			{
				Name:         "wlan1",
				HardwareAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			},
		},
	}

	service := wifisvc.New(h)

	got, err := service.GetAddresses()
	require.NoError(t, err)
	require.Equal(t, []net.HardwareAddr{
		{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
	}, got)
	require.Equal(t, 1, h.calls)
}

func TestWiFiService_GetAddresses_InterfacesError(t *testing.T) {
	t.Parallel()

	h := &handleMock{err: errBoom}
	service := wifisvc.New(h)

	got, err := service.GetAddresses()
	require.Nil(t, got)
	require.Error(t, err)
	require.ErrorContains(t, err, "getting interfaces:")
	require.ErrorIs(t, err, errBoom)
	require.Equal(t, 1, h.calls)
}

func TestWiFiService_GetNames_OK(t *testing.T) {
	t.Parallel()

	h := &handleMock{
		ifaces: []*wifi.Interface{
			{Name: "wlan0"},
			{Name: "wlan1"},
		},
	}

	service := wifisvc.New(h)

	got, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"wlan0", "wlan1"}, got)
	require.Equal(t, 1, h.calls)
}

func TestWiFiService_GetNames_InterfacesError(t *testing.T) {
	t.Parallel()

	h := &handleMock{err: errBoom}
	service := wifisvc.New(h)

	got, err := service.GetNames()
	require.Nil(t, got)
	require.Error(t, err)
	require.ErrorContains(t, err, "getting interfaces:")
	require.ErrorIs(t, err, errBoom)
	require.Equal(t, 1, h.calls)
}
