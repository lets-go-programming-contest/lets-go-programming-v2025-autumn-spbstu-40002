package wifi

import (
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/mdlayher/wifi"
)

func TestWiFiService_GetAddresses_OK(t *testing.T) {
	m := &WiFiHandleMock{}

	ifaces := []*wifi.Interface{
		{Name: "wlan0", HardwareAddr: net.HardwareAddr{0, 1, 2, 3, 4, 5}},
		{Name: "wlan1", HardwareAddr: net.HardwareAddr{10, 11, 12, 13, 14, 15}},
	}

	m.On("Interfaces").Return(ifaces, nil).Once()

	service := New(m)
	got, err := service.GetAddresses()
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 addrs, got: %d", len(got))
	}
	if got[0].String() != ifaces[0].HardwareAddr.String() || got[1].String() != ifaces[1].HardwareAddr.String() {
		t.Fatalf("unexpected addrs: %#v", got)
	}

	m.AssertExpectations(t)
}

func TestWiFiService_GetAddresses_Error(t *testing.T) {
	m := &WiFiHandleMock{}
	m.On("Interfaces").Return(nil, errors.New("boom")).Once()

	service := New(m)
	_, err := service.GetAddresses()
	if err == nil || !strings.Contains(err.Error(), "getting interfaces:") {
		t.Fatalf("expected wrapped error, got: %v", err)
	}

	m.AssertExpectations(t)
}

func TestWiFiService_GetNames_OK(t *testing.T) {
	m := &WiFiHandleMock{}

	ifaces := []*wifi.Interface{
		{Name: "wlan0"},
		{Name: "wlan1"},
	}

	m.On("Interfaces").Return(ifaces, nil).Once()

	service := New(m)
	got, err := service.GetNames()
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if strings.Join(got, ",") != "wlan0,wlan1" {
		t.Fatalf("unexpected names: %#v", got)
	}

	m.AssertExpectations(t)
}

func TestWiFiService_GetNames_Error(t *testing.T) {
	m := &WiFiHandleMock{}
	m.On("Interfaces").Return(nil, errors.New("boom")).Once()

	service := New(m)
	_, err := service.GetNames()
	if err == nil || !strings.Contains(err.Error(), "getting interfaces:") {
		t.Fatalf("expected wrapped error, got: %v", err)
	}

	m.AssertExpectations(t)
}
