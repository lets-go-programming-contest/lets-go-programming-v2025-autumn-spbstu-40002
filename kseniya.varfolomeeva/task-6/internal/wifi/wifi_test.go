package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mywifi "task-6/internal/wifi"
)

var errExpected = errors.New("expected error")

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var err error
	if args.Error(1) != nil {
		err = fmt.Errorf("mock error: %w", args.Error(1))
	}

	if args.Get(0) == nil {
		return nil, err
	}

	if ifaces, ok := args.Get(0).([]*wifi.Interface); ok {
		return ifaces, err
	}

	return nil, err
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockHandle := &MockWiFiHandle{}
	service := mywifi.New(mockHandle)
	require.Equal(t, mockHandle, service.WiFi)
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	type testCase struct {
		addrs []string
		err   error
	}

	cases := []testCase{
		{addrs: []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"}},
		{addrs: []string{}},
		{err: errExpected},
	}

	for i, tc := range cases {
		mockHandle := &MockWiFiHandle{}
		service := mywifi.WiFiService{WiFi: mockHandle}

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

		mockHandle.AssertExpectations(t)
	}
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	type testCase struct {
		addrs []string
		err   error
	}

	cases := []testCase{
		{addrs: []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"}},
		{addrs: []string{}},
		{err: errExpected},
	}

	for i, tc := range cases {
		mockHandle := &MockWiFiHandle{}
		service := mywifi.WiFiService{WiFi: mockHandle}

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

		mockHandle.AssertExpectations(t)
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
