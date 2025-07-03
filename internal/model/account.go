package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Account represents an account in the system
type Account struct {
	ID        int64           `json:"account_id" gorm:"column:account_id;primaryKey"`
	Balance   decimal.Decimal `json:"balance" gorm:"column:balance;type:decimal(20,8);not null;default:0"`
	CreatedAt time.Time       `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName returns the table name for GORM
func (Account) TableName() string {
	return "accounts"
}

// BeforeCreate GORM hook called before creating a record
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	// Validation can be added here if needed
	return nil
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
