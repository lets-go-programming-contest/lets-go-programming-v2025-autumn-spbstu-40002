package wifi_test

import (
	"errors"
	"net"
	"testing"

	mywifi "github.com/Danya-byte/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

//go:generate mockery --all --testonly --quiet --outpkg wifi_test --output .

var ErrSome = errors.New("some error")

type testcase struct {
	names  []string
	addrs  []string
	errMsg string
}

var casesWiFi = []testcase{ //nolint:gochecknoglobals
	{
		addrs: []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
		names: []string{"eth1", "eth2"},
	},
	{
		errMsg: "getting interfaces",
	},
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	for _, test := range casesWiFi {
		mockWifi := NewWiFiHandle(t)
		wifiService := mywifi.New(mockWifi)

		if test.errMsg != "" {
			mockWifi.On("Interfaces").Return(helperMockIfaces(t, &test), ErrSome)
		} else {
			mockWifi.On("Interfaces").Return(helperMockIfaces(t, &test), nil)
		}

		actualAddrs, err := wifiService.GetAddresses()

		if test.errMsg != "" {
			require.ErrorIs(t, err, ErrSome)
			require.ErrorContains(t, err, test.errMsg)

			return
		}

		require.NoError(t, err)
		require.Equal(t, helperParseMACs(t, test.addrs), actualAddrs)
	}
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	for _, test := range casesWiFi {
		mockWifi := NewWiFiHandle(t)
		wifiService := mywifi.New(mockWifi)

		if test.errMsg != "" {
			mockWifi.On("Interfaces").Return(helperMockIfaces(t, &test), ErrSome)
		} else {
			mockWifi.On("Interfaces").Return(helperMockIfaces(t, &test), nil)
		}

		actualNames, err := wifiService.GetNames()

		if test.errMsg != "" {
			require.ErrorIs(t, err, ErrSome)
			require.ErrorContains(t, err, test.errMsg)

			return
		}

		require.NoError(t, err)
		require.Equal(t, test.names, actualNames)
	}
}

func helperMockIfaces(t *testing.T, test *testcase) []*wifi.Interface {
	t.Helper()

	require.Equal(t, len(test.addrs), len(test.names))

	interfaces := make([]*wifi.Interface, 0, len(test.addrs))

	for i, addr := range test.addrs {
		hwAddr := helperParseMAC(t, addr)
		if hwAddr == nil {
			continue
		}

		iface := &wifi.Interface{
			Index:        i + 1,
			Name:         test.names[i],
			HardwareAddr: hwAddr,
			PHY:          1,
			Device:       1,
			Type:         wifi.InterfaceTypeAPVLAN,
			Frequency:    0,
		}
		interfaces = append(interfaces, iface)
	}

	return interfaces
}

func helperParseMACs(t *testing.T, macStr []string) []net.HardwareAddr {
	t.Helper()

	addrs := make([]net.HardwareAddr, 0, len(macStr))

	for _, addr := range macStr {
		addrs = append(addrs, helperParseMAC(t, addr))
	}

	return addrs
}

func helperParseMAC(t *testing.T, macStr string) net.HardwareAddr {
	t.Helper()

	hwAddr, err := net.ParseMAC(macStr)

	require.NoError(t, err)

	return hwAddr
}
