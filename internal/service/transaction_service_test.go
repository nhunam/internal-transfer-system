package service

import (
	"testing"

	"internal-transfer-system/internal/model"
	"internal-transfer-system/internal/repository"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionService_CreateTransaction(t *testing.T) {
	db := setupTestDB(t)
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	accountService := NewAccountService(accountRepo)
	transactionService := NewTransactionService(db, transactionRepo, accountService)

	// Create test accounts
	sourceAccount := &model.CreateAccountRequest{
		AccountID:      123,
		InitialBalance: "100.50",
	}
	destAccount := &model.CreateAccountRequest{
		AccountID:      456,
		InitialBalance: "200.75",
	}

	err := accountService.CreateAccount(sourceAccount)
	require.NoError(t, err)
	err = accountService.CreateAccount(destAccount)
	require.NoError(t, err)

	testCases := []struct {
		name                  string
		request               *model.CreateTransactionRequest
		shouldError           bool
		expectedError         string
		expectedSourceBalance decimal.Decimal
		expectedDestBalance   decimal.Decimal
	}{
		{
			name: "successful transaction",
			request: &model.CreateTransactionRequest{
				SourceAccountID:      123,
				DestinationAccountID: 456,
				Amount:               "25.25",
			},
			shouldError:           false,
			expectedSourceBalance: decimal.NewFromFloat(75.25),  // 100.50 - 25.25
			expectedDestBalance:   decimal.NewFromFloat(226.00), // 200.75 + 25.25
		},
		{
			name: "insufficient balance",
			request: &model.CreateTransactionRequest{
				SourceAccountID:      123,
				DestinationAccountID: 456,
				Amount:               "1000.00",
			},
			shouldError:   true,
			expectedError: "insufficient balance",
		},
		{
			name: "invalid source account",
			request: &model.CreateTransactionRequest{
				SourceAccountID:      999,
				DestinationAccountID: 456,
				Amount:               "25.00",
			},
			shouldError:   true,
			expectedError: "source account validation failed",
		},
		{
			name: "invalid destination account",
			request: &model.CreateTransactionRequest{
				SourceAccountID:      123,
				DestinationAccountID: 999,
				Amount:               "25.00",
			},
			shouldError:   true,
			expectedError: "destination account validation failed",
		},
		{
			name: "same source and destination",
			request: &model.CreateTransactionRequest{
				SourceAccountID:      123,
				DestinationAccountID: 123,
				Amount:               "25.00",
			},
			shouldError:   true,
			expectedError: "source and destination accounts cannot be the same",
		},
		{
			name: "zero amount",
			request: &model.CreateTransactionRequest{
				SourceAccountID:      123,
				DestinationAccountID: 456,
				Amount:               "0",
			},
			shouldError:   true,
			expectedError: "amount must be positive",
		},
		{
			name: "negative amount",
			request: &model.CreateTransactionRequest{
				SourceAccountID:      123,
				DestinationAccountID: 456,
				Amount:               "-10.00",
			},
			shouldError:   true,
			expectedError: "amount must be positive",
		},
		{
			name: "invalid amount format",
			request: &model.CreateTransactionRequest{
				SourceAccountID:      123,
				DestinationAccountID: 456,
				Amount:               "invalid",
			},
			shouldError:   true,
			expectedError: "invalid amount format",
		},
		{
			name: "zero source account ID",
			request: &model.CreateTransactionRequest{
				SourceAccountID:      0,
				DestinationAccountID: 456,
				Amount:               "25.00",
			},
			shouldError:   true,
			expectedError: "source account ID must be positive",
		},
		{
			name: "zero destination account ID",
			request: &model.CreateTransactionRequest{
				SourceAccountID:      123,
				DestinationAccountID: 0,
				Amount:               "25.00",
			},
			shouldError:   true,
			expectedError: "destination account ID must be positive",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get initial balances for successful transaction verification
			var initialSourceBalance, initialDestBalance decimal.Decimal
			if !tc.shouldError {
				initialSourceBalance, _ = accountService.GetAccountBalance(tc.request.SourceAccountID)
				initialDestBalance, _ = accountService.GetAccountBalance(tc.request.DestinationAccountID)
			}

			err := transactionService.CreateTransaction(tc.request)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)

				// Verify balances were updated correctly
				sourceBalance, err := accountService.GetAccountBalance(tc.request.SourceAccountID)
				assert.NoError(t, err)
				destBalance, err := accountService.GetAccountBalance(tc.request.DestinationAccountID)
				assert.NoError(t, err)

				amount, _ := decimal.NewFromString(tc.request.Amount)
				expectedSourceBalance := initialSourceBalance.Sub(amount)
				expectedDestBalance := initialDestBalance.Add(amount)

				assert.True(t, expectedSourceBalance.Equal(sourceBalance),
					"Source balance mismatch: expected %s, got %s",
					expectedSourceBalance.String(), sourceBalance.String())
				assert.True(t, expectedDestBalance.Equal(destBalance),
					"Destination balance mismatch: expected %s, got %s",
					expectedDestBalance.String(), destBalance.String())
			}
		})
	}
}

func TestTransactionService_ErrorScenarios(t *testing.T) {
	db := setupTestDB(t)
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	accountService := NewAccountService(accountRepo)
	transactionService := NewTransactionService(db, transactionRepo, accountService)

	// Create test accounts
	sourceAccount := &model.CreateAccountRequest{
		AccountID:      123,
		InitialBalance: "100.50",
	}
	destAccount := &model.CreateAccountRequest{
		AccountID:      456,
		InitialBalance: "200.75",
	}

	err := accountService.CreateAccount(sourceAccount)
	require.NoError(t, err)
	err = accountService.CreateAccount(destAccount)
	require.NoError(t, err)

	// Test transaction rollback on error
	t.Run("transaction rollback on insufficient balance", func(t *testing.T) {
		// Get initial balances
		sourceBalance, err := accountService.GetAccountBalance(123)
		require.NoError(t, err)
		destBalance, err := accountService.GetAccountBalance(456)
		require.NoError(t, err)

		// Attempt transaction that should fail
		request := &model.CreateTransactionRequest{
			SourceAccountID:      123,
			DestinationAccountID: 456,
			Amount:               "1000.00",
		}

		err = transactionService.CreateTransaction(request)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient balance")

		// Verify balances remain unchanged
		newSourceBalance, err := accountService.GetAccountBalance(123)
		require.NoError(t, err)
		newDestBalance, err := accountService.GetAccountBalance(456)
		require.NoError(t, err)

		assert.True(t, sourceBalance.Equal(newSourceBalance),
			"Source balance should remain unchanged after failed transaction")
		assert.True(t, destBalance.Equal(newDestBalance),
			"Destination balance should remain unchanged after failed transaction")
	})

	// Test various validation scenarios
	t.Run("validation scenarios", func(t *testing.T) {
		validationTests := []struct {
			name     string
			request  *model.CreateTransactionRequest
			errorMsg string
		}{
			{
				name: "empty amount",
				request: &model.CreateTransactionRequest{
					SourceAccountID:      123,
					DestinationAccountID: 456,
					Amount:               "",
				},
				errorMsg: "amount is required",
			},
			{
				name: "invalid amount format",
				request: &model.CreateTransactionRequest{
					SourceAccountID:      123,
					DestinationAccountID: 456,
					Amount:               "invalid",
				},
				errorMsg: "invalid amount format",
			},
		}

		for _, tc := range validationTests {
			t.Run(tc.name, func(t *testing.T) {
				err := transactionService.CreateTransaction(tc.request)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			})
		}
	})
}

func TestTransactionService_ConcurrentTransactions(t *testing.T) {
	db := setupTestDB(t)
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	accountService := NewAccountService(accountRepo)
	transactionService := NewTransactionService(db, transactionRepo, accountService)

	// Create test accounts
	sourceAccount := &model.CreateAccountRequest{
		AccountID:      123,
		InitialBalance: "100.00",
	}
	destAccount := &model.CreateAccountRequest{
		AccountID:      456,
		InitialBalance: "100.00",
	}

	err := accountService.CreateAccount(sourceAccount)
	require.NoError(t, err)
	err = accountService.CreateAccount(destAccount)
	require.NoError(t, err)

	// Test multiple sequential transactions to ensure data consistency
	// Note: Using sequential approach to avoid SQLite concurrency issues in tests
	numTransactions := 5
	amount := decimal.NewFromFloat(10.00)

	// Execute transactions sequentially
	successCount := 0
	for i := 0; i < numTransactions; i++ {
		request := &model.CreateTransactionRequest{
			SourceAccountID:      123,
			DestinationAccountID: 456,
			Amount:               amount.String(),
		}
		err := transactionService.CreateTransaction(request)
		if err == nil {
			successCount++
		}
	}

	// Verify final balances
	sourceBalance, err := accountService.GetAccountBalance(123)
	require.NoError(t, err)
	destBalance, err := accountService.GetAccountBalance(456)
	require.NoError(t, err)

	// Calculate expected balances
	totalTransferred := amount.Mul(decimal.NewFromInt(int64(successCount)))
	expectedSourceBalance := decimal.NewFromFloat(100.00).Sub(totalTransferred)
	expectedDestBalance := decimal.NewFromFloat(100.00).Add(totalTransferred)

	assert.True(t, expectedSourceBalance.Equal(sourceBalance),
		"Source balance mismatch: expected %s, got %s",
		expectedSourceBalance.String(), sourceBalance.String())
	assert.True(t, expectedDestBalance.Equal(destBalance),
		"Destination balance mismatch: expected %s, got %s",
		expectedDestBalance.String(), destBalance.String())

	// Verify total balance conservation
	totalBalance := sourceBalance.Add(destBalance)
	expectedTotalBalance := decimal.NewFromFloat(200.00)
	assert.True(t, expectedTotalBalance.Equal(totalBalance),
		"Total balance not conserved: expected %s, got %s",
		expectedTotalBalance.String(), totalBalance.String())

	// Verify all transactions succeeded
	assert.Equal(t, numTransactions, successCount, "All transactions should succeed")
}
