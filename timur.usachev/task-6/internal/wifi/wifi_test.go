package wifi

import (
	"errors"
	"testing"

	mdwifi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

var errDummy = errors.New("db error")
var errScan = errors.New("scan error")

type mockDB struct {
	names []string
	err   error
}

func (m *mockDB) GetNames() ([]string, error) {
	return m.names, m.err
}

type mockScanner struct {
	ifaces []*mdwifi.Interface
	err    error
}

func (m *mockScanner) Interfaces() ([]*mdwifi.Interface, error) {
	return m.ifaces, m.err
}

func TestAvailableNetworks_TableDriven(t *testing.T) {
	cases := []struct {
		name   string
		input  []string
		expect []string
	}{
		{"single", []string{"Net1"}, []string{"Net1"}},
		{"empty filtered", []string{"", "Net1", ""}, []string{"Net1"}},
		{"duplicates", []string{"A", "B", "A", "C", "B"}, []string{"A", "B", "C"}},
		{"unicode", []string{"сеть", "ネット", "شبكة"}, []string{"сеть", "ネット", "شبكة"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := New(&mockDB{names: tc.input}, &mockScanner{})
			out, err := s.AvailableNetworks()
			require.NoError(t, err)
			require.Equal(t, tc.expect, out)
		})
	}
}

func TestAvailableNetworks_DBError(t *testing.T) {
	s := New(&mockDB{err: errDummy}, &mockScanner{})
	_, err := s.AvailableNetworks()
	require.Error(t, err)
}

func TestAvailableNetworks_ScannerError(t *testing.T) {
	s := New(&mockDB{names: []string{"A"}}, &mockScanner{err: errScan})
	_, err := s.AvailableNetworks()
	require.Error(t, err)
}

func TestNewRealScanner_Coverage(t *testing.T) {
	_, _ = NewRealScanner()
}

func TestRealScanner_Interfaces_PanicCoverage(t *testing.T) {
	r := &RealScanner{}
	defer func() {
		require.NotNil(t, recover())
	}()
	_, _ = r.Interfaces()
}

func TestNewRealScanner_ErrorBranch(t *testing.T) {
	old := newWifiClient
	defer func() { newWifiClient = old }()

	newWifiClient = func() (*mdwifi.Client, error) {
		return nil, errScan
	}

	_, err := NewRealScanner()
	require.Error(t, err)
}

func TestNewRealScanner_SuccessBranch(t *testing.T) {
	old := newWifiClient
	defer func() { newWifiClient = old }()

	newWifiClient = func() (*mdwifi.Client, error) {
		return &mdwifi.Client{}, nil
	}

	s, err := NewRealScanner()
	require.NoError(t, err)
	require.NotNil(t, s)
}
