package repository

import (
	"fmt"

	"internal-transfer-system/internal/model"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// AccountRepository handles database operations for accounts
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// Create creates a new account in the database
func (r *AccountRepository) Create(accountID int64, initialBalance decimal.Decimal) error {
	account := &model.Account{
		ID:      accountID,
		Balance: initialBalance,
	}

	if err := r.db.Create(account).Error; err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}

// GetByID retrieves an account by its ID
func (r *AccountRepository) GetByID(accountID int64) (*model.Account, error) {
	var account model.Account

	if err := r.db.Where("account_id = ?", accountID).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

// UpdateBalance updates the account balance
func (r *AccountRepository) UpdateBalance(accountID int64, newBalance decimal.Decimal) error {
	result := r.db.Model(&model.Account{}).Where("account_id = ?", accountID).Update("balance", newBalance)

	if result.Error != nil {
		return fmt.Errorf("failed to update account balance: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

// Exists checks if an account exists
func (r *AccountRepository) Exists(accountID int64) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Account{}).Where("account_id = ?", accountID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check account existence: %w", err)
	}

	return count > 0, nil
}
