package department

import (
	"fmt"
	"strings"
)

// Manager manages temperature constraints for a department
type Manager struct {
	minTemp int
	maxTemp int
}

// NewManager creates a new temperature manager
func NewManager() *Manager {
	return &Manager{
		minTemp: 15,
		maxTemp: 30,
	}
}

// AddConstraint adds a temperature constraint and returns optimal temperature
func (m *Manager) AddConstraint(constraint string) int {
	// Parse constraint (e.g., ">=30" or "<=25")
	if strings.HasPrefix(constraint, ">=") {
		var temp int
		fmt.Sscanf(constraint, ">=%d", &temp)
		if temp > m.minTemp {
			m.minTemp = temp
		}
	} else if strings.HasPrefix(constraint, "<=") {
		var temp int
		fmt.Sscanf(constraint, "<=%d", &temp)
		if temp < m.maxTemp {
			m.maxTemp = temp
		}
	}

	// Check if valid range exists
	if m.minTemp > m.maxTemp {
		return -1
	}

	// Return the minimum valid temperature
	return m.minTemp
}