package repository

import (
	"database/sql"
	"fmt"
	"time"

	"internal-transfer-system/internal/model"

	"github.com/shopspring/decimal"
)

// AccountRepository handles database operations for accounts
type AccountRepository struct {
	db *sql.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// Create creates a new account in the database
func (r *AccountRepository) Create(accountID int64, initialBalance decimal.Decimal) error {
	query := `
		INSERT INTO accounts (account_id, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`

	now := time.Now()
	_, err := r.db.Exec(query, accountID, initialBalance, now, now)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}

// GetByID retrieves an account by its ID
func (r *AccountRepository) GetByID(accountID int64) (*model.Account, error) {
	query := `
		SELECT account_id, balance, created_at, updated_at
		FROM accounts
		WHERE account_id = $1
	`

	var account model.Account
	err := r.db.QueryRow(query, accountID).Scan(
		&account.ID,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

// UpdateBalance updates the account balance
func (r *AccountRepository) UpdateBalance(accountID int64, newBalance decimal.Decimal) error {
	query := `
		UPDATE accounts
		SET balance = $1, updated_at = $2
		WHERE account_id = $3
	`

	result, err := r.db.Exec(query, newBalance, time.Now(), accountID)
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

// Exists checks if an account exists
func (r *AccountRepository) Exists(accountID int64) (bool, error) {
	query := `SELECT 1 FROM accounts WHERE account_id = $1`

	var exists int
	err := r.db.QueryRow(query, accountID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check account existence: %w", err)
	}

	return true, nil
}
