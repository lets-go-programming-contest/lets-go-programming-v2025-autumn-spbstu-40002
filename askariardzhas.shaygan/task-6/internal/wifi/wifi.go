package wifi

import (
	"fmt"
	"net"

	"github.com/mdlayher/wifi"
)

type WiFiInterface interface {
	Interfaces() ([]*wifi.Interface, error)
}

type WiFiManager struct {
	Handler WiFiInterface
}

func CreateManager(handler WiFiInterface) WiFiManager {
	return WiFiManager{Handler: handler}
}

func (manager WiFiManager) GetMACAddresses() ([]net.HardwareAddr, error) {
	ifaces, err := manager.Handler.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve interfaces: %w", err)
	}

	addresses := make([]net.HardwareAddr, 0, len(ifaces))

	for _, iface := range ifaces {
		addresses = append(addresses, iface.HardwareAddr)
	}

	return addresses, nil
}

func (manager WiFiManager) GetInterfaceNames() ([]string, error) {
	ifaces, err := manager.Handler.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve interfaces: %w", err)
	}

	names := make([]string, 0, len(ifaces))

	for _, iface := range ifaces {
		names = append(names, iface.Name)
	}

	return names, nil
}
