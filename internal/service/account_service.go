package service

import (
	"fmt"

	"internal-transfer-system/internal/model"
	"internal-transfer-system/internal/repository"

	"github.com/shopspring/decimal"
)

// AccountService handles business logic for accounts
type AccountService struct {
	accountRepo *repository.AccountRepository
}

// NewAccountService creates a new account service
func NewAccountService(accountRepo *repository.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}

// CreateAccount creates a new account with initial balance
func (s *AccountService) CreateAccount(request *model.CreateAccountRequest) error {
	// Validate account ID
	if request.AccountID <= 0 {
		return fmt.Errorf("account ID must be positive")
	}

	// Check if account already exists
	exists, err := s.accountRepo.Exists(request.AccountID)
	if err != nil {
		return fmt.Errorf("failed to check account existence: %w", err)
	}
	if exists {
		return fmt.Errorf("account already exists")
	}

	// Parse initial balance
	initialBalance, err := decimal.NewFromString(request.InitialBalance)
	if err != nil {
		return fmt.Errorf("invalid initial balance format: %w", err)
	}

	// Validate initial balance (must be non-negative)
	if initialBalance.IsNegative() {
		return fmt.Errorf("initial balance cannot be negative")
	}

	// Create account
	if err := s.accountRepo.Create(request.AccountID, initialBalance); err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}

// GetAccount retrieves an account by ID
func (s *AccountService) GetAccount(accountID int64) (*model.AccountResponse, error) {
	if accountID <= 0 {
		return nil, fmt.Errorf("account ID must be positive")
	}

	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &model.AccountResponse{
		AccountID: account.ID,
		Balance:   account.Balance.String(),
	}, nil
}

// ValidateAccount checks if an account exists and is valid for transactions
func (s *AccountService) ValidateAccount(accountID int64) error {
	if accountID <= 0 {
		return fmt.Errorf("account ID must be positive")
	}

	exists, err := s.accountRepo.Exists(accountID)
	if err != nil {
		return fmt.Errorf("failed to check account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("account does not exist")
	}

	return nil
}

// GetAccountBalance retrieves the current balance of an account
func (s *AccountService) GetAccountBalance(accountID int64) (decimal.Decimal, error) {
	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get account: %w", err)
	}

	return account.Balance, nil
}

// UpdateAccountBalance updates the account balance
func (s *AccountService) UpdateAccountBalance(accountID int64, newBalance decimal.Decimal) error {
	if newBalance.IsNegative() {
		return fmt.Errorf("account balance cannot be negative")
	}

	if err := s.accountRepo.UpdateBalance(accountID, newBalance); err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	return nil
}
