package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Account represents an account in the system
type Account struct {
	ID        int64           `json:"account_id" db:"account_id"`
	Balance   decimal.Decimal `json:"balance" db:"balance"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

// CreateAccountRequest represents the request payload for creating an account
type CreateAccountRequest struct {
	AccountID      int64  `json:"account_id" binding:"required"`
	InitialBalance string `json:"initial_balance" binding:"required"`
}

// AccountResponse represents the response for account queries
type AccountResponse struct {
	AccountID int64  `json:"account_id"`
	Balance   string `json:"balance"`
}
