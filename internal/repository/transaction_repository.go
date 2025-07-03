package repository

import (
	"database/sql"
	"fmt"
	"time"

	"internal-transfer-system/internal/model"

	"github.com/shopspring/decimal"
)

// TransactionRepository handles database operations for transactions
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create creates a new transaction in the database
func (r *TransactionRepository) Create(sourceAccountID, destinationAccountID int64, amount decimal.Decimal) (*model.Transaction, error) {
	query := `
		INSERT INTO transactions (source_account_id, destination_account_id, amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING transaction_id, source_account_id, destination_account_id, amount, status, created_at, updated_at
	`

	now := time.Now()
	var transaction model.Transaction

	err := r.db.QueryRow(query, sourceAccountID, destinationAccountID, amount, model.TransactionStatusPending, now, now).Scan(
		&transaction.ID,
		&transaction.SourceAccountID,
		&transaction.DestinationAccountID,
		&transaction.Amount,
		&transaction.Status,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return &transaction, nil
}

// UpdateStatus updates the transaction status
func (r *TransactionRepository) UpdateStatus(transactionID int64, status string) error {
	query := `
		UPDATE transactions
		SET status = $1, updated_at = $2
		WHERE transaction_id = $3
	`

	result, err := r.db.Exec(query, status, time.Now(), transactionID)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}

// GetByID retrieves a transaction by its ID
func (r *TransactionRepository) GetByID(transactionID int64) (*model.Transaction, error) {
	query := `
		SELECT transaction_id, source_account_id, destination_account_id, amount, status, created_at, updated_at
		FROM transactions
		WHERE transaction_id = $1
	`

	var transaction model.Transaction
	err := r.db.QueryRow(query, transactionID).Scan(
		&transaction.ID,
		&transaction.SourceAccountID,
		&transaction.DestinationAccountID,
		&transaction.Amount,
		&transaction.Status,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &transaction, nil
}

// GetByAccountID retrieves transactions for a specific account
func (r *TransactionRepository) GetByAccountID(accountID int64, limit, offset int) ([]model.Transaction, error) {
	query := `
		SELECT transaction_id, source_account_id, destination_account_id, amount, status, created_at, updated_at
		FROM transactions
		WHERE source_account_id = $1 OR destination_account_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.SourceAccountID,
			&transaction.DestinationAccountID,
			&transaction.Amount,
			&transaction.Status,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
