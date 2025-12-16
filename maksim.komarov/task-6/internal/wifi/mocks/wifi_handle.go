package mocks

import (
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type WiFiHandle struct {
	mock.Mock
}

func (m *WiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	var r0 []*wifi.Interface
	if v := args.Get(0); v != nil {
		r0 = v.([]*wifi.Interface)
	}
	return r0, args.Error(1)
}

func NewWiFiHandle(t interface {
	mock.TestingT
	Cleanup(func())
}) *WiFiHandle {
	m := &WiFiHandle{}
	m.Mock.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}
