package wifi_test

import (
	wifi "github.com/mdlayher/wifi"
	mock "github.com/stretchr/testify/mock"
)

type WiFiHandle struct {
	mock.Mock
}

func (_m *WiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Interfaces")
	}

	var (
		r0 []*wifi.Interface
		r1 error
	)

	if rf, ok := ret.Get(0).(func() ([]*wifi.Interface, error)); ok {
		return rf()
	}

	if rf, ok := ret.Get(0).(func() []*wifi.Interface); ok {
		r0 = rf()
	} else if v := ret.Get(0); v != nil {
		val, ok := v.([]*wifi.Interface)
		if ok {
			r0 = val
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1 //nolint:wrapcheck
}

func NewWiFiHandle(t interface {
	mock.TestingT
	Cleanup(f func())
},
) *WiFiHandle {
	h := &WiFiHandle{}
	h.Mock.Test(t)

	t.Cleanup(func() { h.AssertExpectations(t) })

	return h
}
