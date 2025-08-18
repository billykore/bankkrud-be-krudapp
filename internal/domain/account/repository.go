package account

// Repository defines a contract for account data access and persistence operations.
type Repository interface {
	// Get retrieves an account from the repository by its account number.
	Get(accountNumber string) (Account, error)
}
