package repository

import (
	"fmt"

	"internal-transfer-system/internal/model"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TransactionRepository handles database operations for transactions
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create creates a new transaction in the database
func (r *TransactionRepository) Create(sourceAccountID, destinationAccountID int64, amount decimal.Decimal) (*model.Transaction, error) {
	transaction := &model.Transaction{
		SourceAccountID:      sourceAccountID,
		DestinationAccountID: destinationAccountID,
		Amount:               amount,
		Status:               model.TransactionStatusPending,
	}

	if err := r.db.Create(transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}

// UpdateStatus updates the transaction status
func (r *TransactionRepository) UpdateStatus(transactionID int64, status string) error {
	result := r.db.Model(&model.Transaction{}).Where("transaction_id = ?", transactionID).Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update transaction status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}

// GetByID retrieves a transaction by its ID
func (r *TransactionRepository) GetByID(transactionID int64) (*model.Transaction, error) {
	var transaction model.Transaction

	if err := r.db.Where("transaction_id = ?", transactionID).First(&transaction).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &transaction, nil
}

// GetByAccountID retrieves transactions for a specific account
func (r *TransactionRepository) GetByAccountID(accountID int64, limit, offset int) ([]model.Transaction, error) {
	var transactions []model.Transaction

	if err := r.db.Where("source_account_id = ? OR destination_account_id = ?", accountID, accountID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}
