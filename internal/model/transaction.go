package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Transaction represents a financial transaction between accounts
type Transaction struct {
	ID                   int64           `json:"transaction_id" db:"transaction_id"`
	SourceAccountID      int64           `json:"source_account_id" db:"source_account_id"`
	DestinationAccountID int64           `json:"destination_account_id" db:"destination_account_id"`
	Amount               decimal.Decimal `json:"amount" db:"amount"`
	Status               string          `json:"status" db:"status"`
	CreatedAt            time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at" db:"updated_at"`
}

// CreateTransactionRequest represents the request payload for creating a transaction
type CreateTransactionRequest struct {
	SourceAccountID      int64  `json:"source_account_id" binding:"required"`
	DestinationAccountID int64  `json:"destination_account_id" binding:"required"`
	Amount               string `json:"amount" binding:"required"`
}

// TransactionStatus constants
const (
	TransactionStatusPending   = "pending"
	TransactionStatusCompleted = "completed"
	TransactionStatusFailed    = "failed"
)
