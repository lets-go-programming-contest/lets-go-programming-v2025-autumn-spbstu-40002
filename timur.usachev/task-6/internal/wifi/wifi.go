package wifi

import mdwifi "github.com/mdlayher/wifi"

type DB interface {
	GetNames() ([]string, error)
}

type Scanner interface {
	Interfaces() ([]*mdwifi.Interface, error)
}

var newWifiClient = mdwifi.New

type RealScanner struct {
	client *mdwifi.Client
}

func NewRealScanner() (*RealScanner, error) {
	c, err := newWifiClient()
	if err != nil {
		return nil, err
	}
	return &RealScanner{client: c}, nil
}

func (r *RealScanner) Interfaces() ([]*mdwifi.Interface, error) {
	return r.client.Interfaces()
}

type Service struct {
	DB      DB
	Scanner Scanner
}

func New(db DB, scanner Scanner) *Service {
	return &Service{DB: db, Scanner: scanner}
}

func (s *Service) AvailableNetworks() ([]string, error) {
	if _, err := s.Scanner.Interfaces(); err != nil {
		return nil, err
	}
	names, err := s.DB.GetNames()
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{})
	var out []string
	for _, n := range names {
		if n == "" {
			continue
		}
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		out = append(out, n)
	}
	return out, nil
}
