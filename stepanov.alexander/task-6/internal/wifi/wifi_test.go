package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	mywifi "github.com/stepanov.alexander/task-6/internal/wifi"
	"github.com/stretchr/testify/require"
)

type fakeWiFiHandle struct {
	interfaces []*wifi.Interface
	err        error
}

func (f *fakeWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	return f.interfaces, f.err
}

type rowTestSysInfo struct {
	addrs       []string
	names       []string
	errExpected error
}

var testTable = []rowTestSysInfo{ //nolint:gochecknoglobals
	{
		addrs: []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
		names: []string{"eth1", "eth2"},
	},
	{
		addrs:       nil,
		errExpected: ErrExpected,
	},
	{
		addrs: []string{},
		names: []string{},
	},
}

var ErrExpected = errors.New("expected error") //nolint:gochecknoglobals

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	for i, row := range testTable {
		t.Run(fmt.Sprintf("row_%d", i), func(t *testing.T) {
			t.Parallel()

			ifaceSlice := helperMockIfaces(t, row.addrs)
			fake := &fakeWiFiHandle{
				interfaces: ifaceSlice,
				err:        row.errExpected,
			}

			wifiService := mywifi.WiFiService{WiFi: fake}

			actualAddrs, err := wifiService.GetAddresses()
			if row.errExpected != nil {
				require.ErrorIs(t, err, row.errExpected)
				return
			}

			require.NoError(t, err)
			require.Equal(t, helperParseMACs(t, row.addrs), actualAddrs)
		})
	}
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	for i, row := range testTable {
		t.Run(fmt.Sprintf("row_%d", i), func(t *testing.T) {
			t.Parallel()

			ifaceSlice := helperMockIfaces(t, row.addrs)
			fake := &fakeWiFiHandle{
				interfaces: ifaceSlice,
				err:        row.errExpected,
			}

			wifiService := mywifi.WiFiService{WiFi: fake}

			actualNames, err := wifiService.GetNames()
			if row.errExpected != nil {
				require.ErrorIs(t, err, row.errExpected)
				return
			}

			require.NoError(t, err)
			require.Equal(t, row.names, actualNames)
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	fake := &fakeWiFiHandle{}
	wifiService := mywifi.New(fake)

	require.Equal(t, fake, wifiService.WiFi)
}

func helperMockIfaces(t *testing.T, addrs []string) []*wifi.Interface {
	t.Helper()

	interfaces := make([]*wifi.Interface, 0, len(addrs))

	for i, addrStr := range addrs {
		hwAddr := parseMAC(addrStr)
		if hwAddr == nil {
			continue
		}

		iface := &wifi.Interface{
			Index:        i + 1,
			Name:         fmt.Sprintf("eth%d", i+1),
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
		addrs = append(addrs, parseMAC(addr))
	}

	return addrs
}

func parseMAC(macStr string) net.HardwareAddr {
	hwAddr, err := net.ParseMAC(macStr)
	if err != nil {
		return nil
	}

	return hwAddr
}
