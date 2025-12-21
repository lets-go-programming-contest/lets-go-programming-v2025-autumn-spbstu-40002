package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"

	wifiPkg "github.com/Svyatoslav2324/task-6/internal/wifi"
)

var ErrSomeError = errors.New("some error")

type mockWifi struct {
	interfaces []*wifi.Interface
	err        error
}

func (m *mockWifi) Interfaces() ([]*wifi.Interface, error) {
	return m.interfaces, m.err
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockWifi := &mockWifi{}

	service := wifiPkg.New(mockWifi)

	require.Equal(t, mockWifi, service.WiFi)
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		interfaces  []*wifi.Interface
		err         error
		expected    []string
		expectedErr bool
	}{
		{
			name: "success",
			interfaces: []*wifi.Interface{
				{Name: "wifi0"},
				{Name: "wifi1"},
				{Name: "wifi2"},
			},
			expected: []string{"wifi0", "wifi1", "wifi2"},
		},
		{
			name:        "interfaces error",
			err:         ErrSomeError,
			expectedErr: true,
		},
		{
			name:       "no interfaces",
			interfaces: []*wifi.Interface{},
			expected:   []string{},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mockWifi := &mockWifi{
				interfaces: testCase.interfaces,
				err:        testCase.err,
			}

			service := wifiPkg.New(mockWifi)

			names, err := service.GetNames()

			if testCase.expectedErr {
				require.Error(t, err)
				require.Nil(t, names)

				return
			}

			require.NoError(t, err)
			require.Equal(t, testCase.expected, names)
		})
	}
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		interfaces  []*wifi.Interface
		err         error
		expected    []net.HardwareAddr
		expectedErr bool
	}{
		{
			name: "success",
			interfaces: []*wifi.Interface{
				{Name: "wifi0", HardwareAddr: net.HardwareAddr{0x1, 0x2, 0x3, 0x4, 0x5, 0x6}},
				{Name: "wifi1", HardwareAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc}},
				{Name: "wifi2", HardwareAddr: net.HardwareAddr{0xdd, 0xee, 0xff}},
			},
			expected: []net.HardwareAddr{
				{0x1, 0x2, 0x3, 0x4, 0x5, 0x6},
				{0xaa, 0xbb, 0xcc},
				{0xdd, 0xee, 0xff},
			},
		},
		{
			name:        "interfaces error",
			err:         ErrSomeError,
			expectedErr: true,
		},
		{
			name:       "no interfaces",
			interfaces: []*wifi.Interface{},
			expected:   []net.HardwareAddr{},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mockWifi := &mockWifi{
				interfaces: testCase.interfaces,
				err:        testCase.err,
			}

			service := wifiPkg.New(mockWifi)

			addresses, err := service.GetAddresses()

			if testCase.expectedErr {
				require.Error(t, err)
				require.Nil(t, addresses)

				return
			}

			require.NoError(t, err)
			require.Equal(t, testCase.expected, addresses)
		})
	}
}
