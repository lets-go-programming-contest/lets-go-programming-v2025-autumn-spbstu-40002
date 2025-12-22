package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
	myWifi "platon.karyakin/task-6/internal/wifi"
)

//go:generate mockery --all --testonly --quiet --outpkg wifi_test --output .

var errExpected = errors.New("expected error")

type testCase struct {
	addrs []string
	err   error
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockHandle := NewWiFiHandle(t)
	service := myWifi.New(mockHandle)
	require.Equal(t, mockHandle, service.WiFi)
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	cases := []testCase{
		{addrs: []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"}},
		{addrs: []string{}},
		{err: errExpected},
	}

	mockHandle := NewWiFiHandle(t)
	service := myWifi.WiFiService{WiFi: mockHandle}

	for i, tc := range cases {
		mockHandle.On("Interfaces").Unset()
		mockHandle.On("Interfaces").Return(makeIfaces(t, tc.addrs), tc.err)

		got, err := service.GetAddresses()

		if tc.err != nil {
			require.ErrorIs(t, err, tc.err, "case %d", i)
			require.ErrorContains(t, err, "getting interfaces", "case %d", i)
			require.Nil(t, got, "case %d", i)

			continue
		}

		require.NoError(t, err, "case %d", i)
		require.Equal(t, parseMACs(t, tc.addrs), got, "case %d", i)
	}
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	cases := []testCase{
		{addrs: []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"}},
		{addrs: []string{}},
		{err: errExpected},
	}

	mockHandle := NewWiFiHandle(t)
	service := myWifi.WiFiService{WiFi: mockHandle}

	for i, tc := range cases {
		mockHandle.On("Interfaces").Unset()
		mockHandle.On("Interfaces").Return(makeIfaces(t, tc.addrs), tc.err)

		got, err := service.GetNames()

		if tc.err != nil {
			require.ErrorIs(t, err, tc.err, "case %d", i)
			require.ErrorContains(t, err, "getting interfaces", "case %d", i)
			require.Nil(t, got, "case %d", i)

			continue
		}

		require.NoError(t, err, "case %d", i)
		require.Equal(t, wantNames(tc.addrs), got, "case %d", i)
	}
}

func wantNames(addrs []string) []string {
	names := make([]string, 0, len(addrs))
	for i := range addrs {
		names = append(names, fmt.Sprintf("wlan%d", i+1))
	}

	return names
}

func makeIfaces(t *testing.T, addrs []string) []*wifi.Interface {
	t.Helper()

	ifaces := make([]*wifi.Interface, 0, len(addrs))

	for i, macStr := range addrs {
		hw := parseMAC(t, macStr)

		ifaces = append(ifaces, &wifi.Interface{
			Index:        i + 1,
			Name:         fmt.Sprintf("wlan%d", i+1),
			HardwareAddr: hw,
			PHY:          1,
			Device:       1,
			Type:         wifi.InterfaceTypeAPVLAN,
			Frequency:    0,
		})
	}

	return ifaces
}

func parseMACs(t *testing.T, addrs []string) []net.HardwareAddr {
	t.Helper()

	result := make([]net.HardwareAddr, 0, len(addrs))

	for _, s := range addrs {
		result = append(result, parseMAC(t, s))
	}

	return result
}

func parseMAC(t *testing.T, macStr string) net.HardwareAddr {
	t.Helper()

	hw, err := net.ParseMAC(macStr)
	require.NoError(t, err)

	return hw
}
