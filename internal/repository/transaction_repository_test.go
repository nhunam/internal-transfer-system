package repository

import (
	"testing"

	"internal-transfer-system/internal/model"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	transactionRepo := NewTransactionRepository(db)
	accountRepo := NewAccountRepository(db)

	// Create test accounts
	sourceAccountID := int64(123)
	destAccountID := int64(456)
	err := accountRepo.Create(sourceAccountID, decimal.NewFromFloat(100.00))
	require.NoError(t, err)
	err = accountRepo.Create(destAccountID, decimal.NewFromFloat(50.00))
	require.NoError(t, err)

	testCases := []struct {
		name                 string
		sourceAccountID      int64
		destinationAccountID int64
		amount               decimal.Decimal
		shouldError          bool
	}{
		{
			name:                 "successful transaction creation",
			sourceAccountID:      sourceAccountID,
			destinationAccountID: destAccountID,
			amount:               decimal.NewFromFloat(25.50),
			shouldError:          false,
		},
		{
			name:                 "transaction with zero amount",
			sourceAccountID:      sourceAccountID,
			destinationAccountID: destAccountID,
			amount:               decimal.Zero,
			shouldError:          false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			transaction, err := transactionRepo.Create(tc.sourceAccountID, tc.destinationAccountID, tc.amount)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Nil(t, transaction)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, transaction)
				assert.NotZero(t, transaction.ID)
				assert.Equal(t, tc.sourceAccountID, transaction.SourceAccountID)
				assert.Equal(t, tc.destinationAccountID, transaction.DestinationAccountID)
				assert.True(t, tc.amount.Equal(transaction.Amount))
				assert.Equal(t, model.TransactionStatusPending, transaction.Status)
				assert.False(t, transaction.CreatedAt.IsZero())
				assert.False(t, transaction.UpdatedAt.IsZero())
			}
		})
	}
}

func TestTransactionRepository_UpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	transactionRepo := NewTransactionRepository(db)
	accountRepo := NewAccountRepository(db)

	// Create test accounts
	sourceAccountID := int64(123)
	destAccountID := int64(456)
	err := accountRepo.Create(sourceAccountID, decimal.NewFromFloat(100.00))
	require.NoError(t, err)
	err = accountRepo.Create(destAccountID, decimal.NewFromFloat(50.00))
	require.NoError(t, err)

	// Create test transaction
	transaction, err := transactionRepo.Create(sourceAccountID, destAccountID, decimal.NewFromFloat(25.50))
	require.NoError(t, err)

	testCases := []struct {
		name          string
		transactionID int64
		newStatus     string
		shouldError   bool
		errorMsg      string
	}{
		{
			name:          "update to completed status",
			transactionID: transaction.ID,
			newStatus:     model.TransactionStatusCompleted,
			shouldError:   false,
		},
		{
			name:          "update to failed status",
			transactionID: transaction.ID,
			newStatus:     model.TransactionStatusFailed,
			shouldError:   false,
		},
		{
			name:          "non-existent transaction",
			transactionID: 999,
			newStatus:     model.TransactionStatusCompleted,
			shouldError:   true,
			errorMsg:      "transaction not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := transactionRepo.UpdateStatus(tc.transactionID, tc.newStatus)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)

				// Verify the status was updated
				updatedTransaction, err := transactionRepo.GetByID(tc.transactionID)
				assert.NoError(t, err)
				assert.Equal(t, tc.newStatus, updatedTransaction.Status)
			}
		})
	}
}

func TestTransactionRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	transactionRepo := NewTransactionRepository(db)
	accountRepo := NewAccountRepository(db)

	// Create test accounts
	sourceAccountID := int64(123)
	destAccountID := int64(456)
	err := accountRepo.Create(sourceAccountID, decimal.NewFromFloat(100.00))
	require.NoError(t, err)
	err = accountRepo.Create(destAccountID, decimal.NewFromFloat(50.00))
	require.NoError(t, err)

	// Create test transaction
	createdTransaction, err := transactionRepo.Create(sourceAccountID, destAccountID, decimal.NewFromFloat(25.50))
	require.NoError(t, err)

	testCases := []struct {
		name          string
		transactionID int64
		shouldError   bool
		errorMsg      string
	}{
		{
			name:          "existing transaction",
			transactionID: createdTransaction.ID,
			shouldError:   false,
		},
		{
			name:          "non-existent transaction",
			transactionID: 999,
			shouldError:   true,
			errorMsg:      "transaction not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			transaction, err := transactionRepo.GetByID(tc.transactionID)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, transaction)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, transaction)
				assert.Equal(t, tc.transactionID, transaction.ID)
				assert.Equal(t, sourceAccountID, transaction.SourceAccountID)
				assert.Equal(t, destAccountID, transaction.DestinationAccountID)
				assert.True(t, decimal.NewFromFloat(25.50).Equal(transaction.Amount))
			}
		})
	}
}

func TestTransactionRepository_GetByAccountID(t *testing.T) {
	db := setupTestDB(t)
	transactionRepo := NewTransactionRepository(db)
	accountRepo := NewAccountRepository(db)

	// Create test accounts
	account1ID := int64(123)
	account2ID := int64(456)
	account3ID := int64(789)
	err := accountRepo.Create(account1ID, decimal.NewFromFloat(100.00))
	require.NoError(t, err)
	err = accountRepo.Create(account2ID, decimal.NewFromFloat(50.00))
	require.NoError(t, err)
	err = accountRepo.Create(account3ID, decimal.NewFromFloat(75.00))
	require.NoError(t, err)

	// Create test transactions
	_, err = transactionRepo.Create(account1ID, account2ID, decimal.NewFromFloat(25.00))
	require.NoError(t, err)
	_, err = transactionRepo.Create(account2ID, account1ID, decimal.NewFromFloat(10.00))
	require.NoError(t, err)
	_, err = transactionRepo.Create(account1ID, account3ID, decimal.NewFromFloat(15.00))
	require.NoError(t, err)
	_, err = transactionRepo.Create(account2ID, account3ID, decimal.NewFromFloat(5.00))
	require.NoError(t, err)

	testCases := []struct {
		name                string
		accountID           int64
		limit               int
		offset              int
		expectedCount       int
		shouldContainAmount []decimal.Decimal
	}{
		{
			name:                "account1 transactions (all)",
			accountID:           account1ID,
			limit:               10,
			offset:              0,
			expectedCount:       3, // account1 appears in 3 transactions
			shouldContainAmount: []decimal.Decimal{decimal.NewFromFloat(25.00), decimal.NewFromFloat(10.00), decimal.NewFromFloat(15.00)},
		},
		{
			name:                "account2 transactions (all)",
			accountID:           account2ID,
			limit:               10,
			offset:              0,
			expectedCount:       3, // account2 appears in 3 transactions
			shouldContainAmount: []decimal.Decimal{decimal.NewFromFloat(25.00), decimal.NewFromFloat(10.00), decimal.NewFromFloat(5.00)},
		},
		{
			name:          "account3 transactions (limited)",
			accountID:     account3ID,
			limit:         1,
			offset:        0,
			expectedCount: 1, // limit to 1 transaction
		},
		{
			name:          "non-existent account",
			accountID:     999,
			limit:         10,
			offset:        0,
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			transactions, err := transactionRepo.GetByAccountID(tc.accountID, tc.limit, tc.offset)
			assert.NoError(t, err)
			assert.Len(t, transactions, tc.expectedCount)

			// Verify transactions are ordered by created_at DESC (most recent first)
			if len(transactions) > 1 {
				for i := 0; i < len(transactions)-1; i++ {
					assert.True(t, transactions[i].CreatedAt.After(transactions[i+1].CreatedAt) ||
						transactions[i].CreatedAt.Equal(transactions[i+1].CreatedAt))
				}
			}

			// Verify transactions involve the specified account
			for _, transaction := range transactions {
				assert.True(t, transaction.SourceAccountID == tc.accountID ||
					transaction.DestinationAccountID == tc.accountID)
			}
		})
	}
}
