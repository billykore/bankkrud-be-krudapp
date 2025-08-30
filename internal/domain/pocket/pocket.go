// Package pocket contains the domain entities for the pocket module.
package pocket

const StatusActive = "active"

// Pocket represents a financial pocket entity.
type Pocket struct {
	ID            uint64
	AccountNumber string
	Name          string
	Status        string
}

func (p *Pocket) NotActive() bool {
	return p.Status != StatusActive
}
