// Package pocket contains the domain entities for the pocket module.
package pocket

// Pocket represents a financial pocket entity.
type Pocket struct {
	ID            int64
	AccountNumber string
	Name          string
	Status        string
}
