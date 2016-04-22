package internalAccount

import (
	"errors"
)

var (
	//ErrNoRowFound - no row was found in store
	ErrNoRowFound = errors.New("No Row Found")
	//ErrMissingAccountID - the userid was not passed
	ErrMissingAccountID = errors.New("Missing AccountID")
)

// Account -
type Account struct {
	AccountID    string        `json:"accountID"`
	AccountItems []AccountItem `json:"accountItems"`
}

//AccountItem -
type AccountItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// AccountStore - storage manager for extra Account Information
type AccountStore interface {
	// Add - adds extra account information
	Add(account Account) error

	// Update - updates extra account information
	Update(account Account) error

	// Delete -
	Delete(accountID string) error

	// GetByAccountID - retrieves Account by userID
	GetByAccountID(accountID string) (*Account, error)
}
