package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Transaction represents a financial transaction between accounts
type Transaction struct {
	ID                   int64           `json:"transaction_id" gorm:"column:transaction_id;primaryKey;autoIncrement"`
	SourceAccountID      int64           `json:"source_account_id" gorm:"column:source_account_id;not null;index"`
	DestinationAccountID int64           `json:"destination_account_id" gorm:"column:destination_account_id;not null;index"`
	Amount               decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(20,8);not null"`
	Status               string          `json:"status" gorm:"column:status;type:varchar(20);not null;default:pending;index"`
	CreatedAt            time.Time       `json:"created_at" gorm:"column:created_at;autoCreateTime;index"`
	UpdatedAt            time.Time       `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`

	// Relations
	SourceAccount      Account `gorm:"foreignKey:SourceAccountID;references:ID"`
	DestinationAccount Account `gorm:"foreignKey:DestinationAccountID;references:ID"`
}

// TableName returns the table name for GORM
func (Transaction) TableName() string {
	return "transactions"
}

// BeforeCreate GORM hook called before creating a record
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	// Set default status if not provided
	if t.Status == "" {
		t.Status = TransactionStatusPending
	}
	return nil
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
