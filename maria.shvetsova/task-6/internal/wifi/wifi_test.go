package wifi_test

import (
	"errors"
	"net"
	"testing"

	wifiext "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ummmsh/task-6/internal/wifi"
)

var (
	errDriverUnavailable = errors.New("driver unavailable")
	errPermissionDenied  = errors.New("permission denied")
)

type manualMockWiFi struct {
	interfacesFunc func() ([]*wifiext.Interface, error)
}

func (m *manualMockWiFi) Interfaces() ([]*wifiext.Interface, error) {
	if m.interfacesFunc != nil {
		return m.interfacesFunc()
	}

	return nil, nil
}

func TestWiFiServiceGetAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		addr1, _ := net.ParseMAC("00:11:22:33:44:55")
		addr2, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")

		mock := &manualMockWiFi{
			interfacesFunc: func() ([]*wifiext.Interface, error) {
				return []*wifiext.Interface{
					{HardwareAddr: addr1},
					{HardwareAddr: addr2},
				}, nil
			},
		}

		service := wifi.New(mock)
		addrs, err := service.GetAddresses()

		require.NoError(t, err)
		assert.Equal(t, []net.HardwareAddr{addr1, addr2}, addrs)
	})

	t.Run("interfaces error", func(t *testing.T) {
		t.Parallel()

		mock := &manualMockWiFi{
			interfacesFunc: func() ([]*wifiext.Interface, error) {
				return nil, errDriverUnavailable
			},
		}

		service := wifi.New(mock)
		addrs, err := service.GetAddresses()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "getting interfaces")
		assert.Nil(t, addrs)
	})
}

func TestWiFiServiceGetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		addr, _ := net.ParseMAC("00:11:22:33:44:55")

		mock := &manualMockWiFi{
			interfacesFunc: func() ([]*wifiext.Interface, error) {
				return []*wifiext.Interface{
					{Name: "wlan0", HardwareAddr: addr},
					{Name: "wlan1", HardwareAddr: addr},
				}, nil
			},
		}

		service := wifi.New(mock)
		names, err := service.GetNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"wlan0", "wlan1"}, names)
	})

	t.Run("interfaces error", func(t *testing.T) {
		t.Parallel()

		mock := &manualMockWiFi{
			interfacesFunc: func() ([]*wifiext.Interface, error) {
				return nil, errPermissionDenied
			},
		}

		service := wifi.New(mock)
		names, err := service.GetNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "getting interfaces")
		assert.Nil(t, names)
	})
}
