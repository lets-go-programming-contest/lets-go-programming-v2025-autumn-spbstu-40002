package wifi

import (
	"errors"
	"net"
	"testing"
)

type MockWiFiHandle struct {
	interfaces []*struct {
		Name         string
		HardwareAddr net.HardwareAddr
	}
	err error
}

func (m *MockWiFiHandle) Interfaces() ([]*struct {
	Name         string
	HardwareAddr net.HardwareAddr
}, error) {
	return m.interfaces, m.err
}

func parseMAC(s string) net.HardwareAddr {
	mac, err := net.ParseMAC(s)
	if err != nil {
		return nil
	}
	return mac
}

func TestWiFiService_GetAddresses(t *testing.T) {
	mac1 := parseMAC("00:11:22:33:44:55")
	mac2 := parseMAC("aa:bb:cc:dd:ee:ff")

	mock := &MockWiFiHandle{
		interfaces: []*struct {
			Name         string
			HardwareAddr net.HardwareAddr
		}{
			{Name: "wlan0", HardwareAddr: mac1},
			{Name: "wlan1", HardwareAddr: mac2},
		},
	}

	service := New(mock)
	addrs, err := service.GetAddresses()
	if err != nil {
		t.Errorf("GetAddresses() error = %v", err)
	}
	if len(addrs) != 2 {
		t.Errorf("GetAddresses() returned %d addresses, want 2", len(addrs))
	}

	mockWithError := &MockWiFiHandle{
		err: errors.New("mock error"),
	}
	service2 := New(mockWithError)

	_, err = service2.GetAddresses()
	if err == nil {
		t.Error("GetAddresses() should return error when Interfaces() fails")
	}

	mockEmpty := &MockWiFiHandle{
		interfaces: []*struct {
			Name         string
			HardwareAddr net.HardwareAddr
		}{},
	}
	service3 := New(mockEmpty)

	addrs, err = service3.GetAddresses()
	if err != nil {
		t.Errorf("GetAddresses() error = %v", err)
	}
	if len(addrs) != 0 {
		t.Errorf("GetAddresses() returned %d addresses, want 0", len(addrs))
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	mac1 := parseMAC("00:11:22:33:44:55")

	mock := &MockWiFiHandle{
		interfaces: []*struct {
			Name         string
			HardwareAddr net.HardwareAddr
		}{
			{Name: "wlan0", HardwareAddr: mac1},
			{Name: "wlan1", HardwareAddr: mac1},
		},
	}

	service := New(mock)
	names, err := service.GetNames()
	if err != nil {
		t.Errorf("GetNames() error = %v", err)
	}
	if len(names) != 2 {
		t.Errorf("GetNames() returned %d names, want 2", len(names))
	}
	if names[0] != "wlan0" || names[1] != "wlan1" {
		t.Errorf("GetNames() returned %v, want [wlan0 wlan1]", names)
	}

	mockWithError := &MockWiFiHandle{
		err: errors.New("mock error"),
	}
	service2 := New(mockWithError)

	_, err = service2.GetNames()
	if err == nil {
		t.Error("GetNames() should return error when Interfaces() fails")
	}

	mockEmpty := &MockWiFiHandle{
		interfaces: []*struct {
			Name         string
			HardwareAddr net.HardwareAddr
		}{},
	}
	service3 := New(mockEmpty)

	names, err = service3.GetNames()
	if err != nil {
		t.Errorf("GetNames() error = %v", err)
	}
	if len(names) != 0 {
		t.Errorf("GetNames() returned %d names, want 0", len(names))
	}
}
