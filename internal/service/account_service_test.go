package service

import (
	"testing"

	"internal-transfer-system/internal/model"
	"internal-transfer-system/internal/repository"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountService_CreateAccount(t *testing.T) {
	db := setupTestDB(t)
	accountRepo := repository.NewAccountRepository(db)
	accountService := NewAccountService(accountRepo)

	testCases := []struct {
		name          string
		request       *model.CreateAccountRequest
		shouldError   bool
		expectedError string
	}{
		{
			name: "successful account creation",
			request: &model.CreateAccountRequest{
				AccountID:      123,
				InitialBalance: "100.50",
			},
			shouldError: false,
		},
		{
			name: "account with zero balance",
			request: &model.CreateAccountRequest{
				AccountID:      456,
				InitialBalance: "0",
			},
			shouldError: false,
		},
		{
			name: "account with large balance",
			request: &model.CreateAccountRequest{
				AccountID:      789,
				InitialBalance: "1000000.12345678",
			},
			shouldError: false,
		},
		{
			name: "duplicate account ID",
			request: &model.CreateAccountRequest{
				AccountID:      123,
				InitialBalance: "200.00",
			},
			shouldError:   true,
			expectedError: "account already exists",
		},
		{
			name: "invalid account ID (zero)",
			request: &model.CreateAccountRequest{
				AccountID:      0,
				InitialBalance: "100.00",
			},
			shouldError:   true,
			expectedError: "account ID must be positive",
		},
		{
			name: "invalid account ID (negative)",
			request: &model.CreateAccountRequest{
				AccountID:      -123,
				InitialBalance: "100.00",
			},
			shouldError:   true,
			expectedError: "account ID must be positive",
		},
		{
			name: "invalid balance format",
			request: &model.CreateAccountRequest{
				AccountID:      999,
				InitialBalance: "invalid",
			},
			shouldError:   true,
			expectedError: "invalid initial balance format",
		},
		{
			name: "negative balance",
			request: &model.CreateAccountRequest{
				AccountID:      999,
				InitialBalance: "-100.00",
			},
			shouldError:   true,
			expectedError: "initial balance cannot be negative",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := accountService.CreateAccount(tc.request)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)

				// Verify the account was created correctly
				response, err := accountService.GetAccount(tc.request.AccountID)
				assert.NoError(t, err)
				assert.Equal(t, tc.request.AccountID, response.AccountID)

				expectedBalance, _ := decimal.NewFromString(tc.request.InitialBalance)
				actualBalance, _ := decimal.NewFromString(response.Balance)
				assert.True(t, expectedBalance.Equal(actualBalance))
			}
		})
	}
}

func TestAccountService_GetAccount(t *testing.T) {
	db := setupTestDB(t)
	accountRepo := repository.NewAccountRepository(db)
	accountService := NewAccountService(accountRepo)

	// Create test account
	createRequest := &model.CreateAccountRequest{
		AccountID:      123,
		InitialBalance: "100.50",
	}
	err := accountService.CreateAccount(createRequest)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		accountID     int64
		shouldError   bool
		expectedError string
	}{
		{
			name:        "existing account",
			accountID:   123,
			shouldError: false,
		},
		{
			name:          "non-existent account",
			accountID:     999,
			shouldError:   true,
			expectedError: "failed to get account",
		},
		{
			name:          "invalid account ID (zero)",
			accountID:     0,
			shouldError:   true,
			expectedError: "account ID must be positive",
		},
		{
			name:          "invalid account ID (negative)",
			accountID:     -123,
			shouldError:   true,
			expectedError: "account ID must be positive",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := accountService.GetAccount(tc.accountID)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tc.accountID, response.AccountID)
				assert.Equal(t, "100.5", response.Balance)
			}
		})
	}
}

func TestAccountService_ValidateAccount(t *testing.T) {
	db := setupTestDB(t)
	accountRepo := repository.NewAccountRepository(db)
	accountService := NewAccountService(accountRepo)

	// Create test account
	createRequest := &model.CreateAccountRequest{
		AccountID:      123,
		InitialBalance: "100.50",
	}
	err := accountService.CreateAccount(createRequest)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		accountID   int64
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "valid existing account",
			accountID:   123,
			shouldError: false,
		},
		{
			name:        "non-existent account",
			accountID:   999,
			shouldError: true,
			errorMsg:    "account does not exist",
		},
		{
			name:        "invalid account ID (zero)",
			accountID:   0,
			shouldError: true,
			errorMsg:    "account ID must be positive",
		},
		{
			name:        "invalid account ID (negative)",
			accountID:   -123,
			shouldError: true,
			errorMsg:    "account ID must be positive",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := accountService.ValidateAccount(tc.accountID)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAccountService_UpdateAccountBalance(t *testing.T) {
	db := setupTestDB(t)
	accountRepo := repository.NewAccountRepository(db)
	accountService := NewAccountService(accountRepo)

	// Create test account
	createRequest := &model.CreateAccountRequest{
		AccountID:      123,
		InitialBalance: "100.50",
	}
	err := accountService.CreateAccount(createRequest)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		accountID   int64
		newBalance  decimal.Decimal
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "successful balance update",
			accountID:   123,
			newBalance:  decimal.NewFromFloat(75.25),
			shouldError: false,
		},
		{
			name:        "update to zero balance",
			accountID:   123,
			newBalance:  decimal.Zero,
			shouldError: false,
		},
		{
			name:        "negative balance should fail",
			accountID:   123,
			newBalance:  decimal.NewFromFloat(-10.00),
			shouldError: true,
			errorMsg:    "account balance cannot be negative",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := accountService.UpdateAccountBalance(tc.accountID, tc.newBalance)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)

				// Verify the balance was updated
				balance, err := accountService.GetAccountBalance(tc.accountID)
				assert.NoError(t, err)
				assert.True(t, tc.newBalance.Equal(balance))
			}
		})
	}
}

func TestAccountService_GetAccountBalance(t *testing.T) {
	db := setupTestDB(t)
	accountRepo := repository.NewAccountRepository(db)
	accountService := NewAccountService(accountRepo)

	// Create test account
	createRequest := &model.CreateAccountRequest{
		AccountID:      123,
		InitialBalance: "100.50",
	}
	err := accountService.CreateAccount(createRequest)
	require.NoError(t, err)

	testCases := []struct {
		name            string
		accountID       int64
		shouldError     bool
		expectedBalance decimal.Decimal
		errorMsg        string
	}{
		{
			name:            "existing account",
			accountID:       123,
			shouldError:     false,
			expectedBalance: decimal.NewFromFloat(100.50),
		},
		{
			name:        "non-existent account",
			accountID:   999,
			shouldError: true,
			errorMsg:    "failed to get account",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			balance, err := accountService.GetAccountBalance(tc.accountID)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.True(t, tc.expectedBalance.Equal(balance))
			}
		})
	}
}
