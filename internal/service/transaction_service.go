package service

import (
	"database/sql"
	"fmt"

	"internal-transfer-system/internal/model"
	"internal-transfer-system/internal/repository"

	"github.com/shopspring/decimal"
)

// TransactionService handles business logic for transactions
type TransactionService struct {
	db              *sql.DB
	transactionRepo *repository.TransactionRepository
	accountService  *AccountService
}

// NewTransactionService creates a new transaction service
func NewTransactionService(db *sql.DB, transactionRepo *repository.TransactionRepository, accountService *AccountService) *TransactionService {
	return &TransactionService{
		db:              db,
		transactionRepo: transactionRepo,
		accountService:  accountService,
	}
}

// CreateTransaction creates and processes a new transaction
func (s *TransactionService) CreateTransaction(request *model.CreateTransactionRequest) error {
	// Validate request
	if err := s.validateTransactionRequest(request); err != nil {
		return err
	}

	// Parse amount
	amount, err := decimal.NewFromString(request.Amount)
	if err != nil {
		return fmt.Errorf("invalid amount format: %w", err)
	}

	// Validate amount
	if amount.IsNegative() || amount.IsZero() {
		return fmt.Errorf("amount must be positive")
	}

	// Validate accounts exist
	if err := s.accountService.ValidateAccount(request.SourceAccountID); err != nil {
		return fmt.Errorf("source account validation failed: %w", err)
	}

	if err := s.accountService.ValidateAccount(request.DestinationAccountID); err != nil {
		return fmt.Errorf("destination account validation failed: %w", err)
	}

	// Process transaction in database transaction
	return s.processTransaction(request.SourceAccountID, request.DestinationAccountID, amount)
}

// validateTransactionRequest validates the transaction request
func (s *TransactionService) validateTransactionRequest(request *model.CreateTransactionRequest) error {
	if request.SourceAccountID <= 0 {
		return fmt.Errorf("source account ID must be positive")
	}

	if request.DestinationAccountID <= 0 {
		return fmt.Errorf("destination account ID must be positive")
	}

	if request.SourceAccountID == request.DestinationAccountID {
		return fmt.Errorf("source and destination accounts cannot be the same")
	}

	if request.Amount == "" {
		return fmt.Errorf("amount is required")
	}

	return nil
}

// processTransaction processes the transaction with proper data integrity
func (s *TransactionService) processTransaction(sourceAccountID, destinationAccountID int64, amount decimal.Decimal) error {
	// Start database transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start database transaction: %w", err)
	}
	defer tx.Rollback()

	// Lock accounts for update to prevent concurrent modifications
	sourceBalance, err := s.getAccountBalanceForUpdate(tx, sourceAccountID)
	if err != nil {
		return fmt.Errorf("failed to get source account balance: %w", err)
	}

	destinationBalance, err := s.getAccountBalanceForUpdate(tx, destinationAccountID)
	if err != nil {
		return fmt.Errorf("failed to get destination account balance: %w", err)
	}

	// Check if source account has sufficient balance
	if sourceBalance.LessThan(amount) {
		return fmt.Errorf("insufficient balance in source account")
	}

	// Calculate new balances
	newSourceBalance := sourceBalance.Sub(amount)
	newDestinationBalance := destinationBalance.Add(amount)

	// Update account balances
	if err := s.updateAccountBalanceInTx(tx, sourceAccountID, newSourceBalance); err != nil {
		return fmt.Errorf("failed to update source account balance: %w", err)
	}

	if err := s.updateAccountBalanceInTx(tx, destinationAccountID, newDestinationBalance); err != nil {
		return fmt.Errorf("failed to update destination account balance: %w", err)
	}

	// Create transaction record
	if err := s.createTransactionInTx(tx, sourceAccountID, destinationAccountID, amount, model.TransactionStatusCompleted); err != nil {
		return fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// getAccountBalanceForUpdate gets account balance with row lock
func (s *TransactionService) getAccountBalanceForUpdate(tx *sql.Tx, accountID int64) (decimal.Decimal, error) {
	query := `SELECT balance FROM accounts WHERE account_id = $1 FOR UPDATE`

	var balance decimal.Decimal
	err := tx.QueryRow(query, accountID).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return decimal.Zero, fmt.Errorf("account not found")
		}
		return decimal.Zero, fmt.Errorf("failed to get account balance: %w", err)
	}

	return balance, nil
}

// updateAccountBalanceInTx updates account balance within a transaction
func (s *TransactionService) updateAccountBalanceInTx(tx *sql.Tx, accountID int64, newBalance decimal.Decimal) error {
	query := `UPDATE accounts SET balance = $1, updated_at = CURRENT_TIMESTAMP WHERE account_id = $2`

	result, err := tx.Exec(query, newBalance, accountID)
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

// createTransactionInTx creates a transaction record within a transaction
func (s *TransactionService) createTransactionInTx(tx *sql.Tx, sourceAccountID, destinationAccountID int64, amount decimal.Decimal, status string) error {
	query := `
		INSERT INTO transactions (source_account_id, destination_account_id, amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err := tx.Exec(query, sourceAccountID, destinationAccountID, amount, status)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}
