package repository

import (
	"testing"

	"internal-transfer-system/internal/model"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the schema
	err = db.AutoMigrate(&model.Account{}, &model.Transaction{})
	require.NoError(t, err)

	return db
}

func TestAccountRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	testCases := []struct {
		name           string
		accountID      int64
		initialBalance decimal.Decimal
		shouldError    bool
	}{
		{
			name:           "successful account creation",
			accountID:      123,
			initialBalance: decimal.NewFromFloat(100.50),
			shouldError:    false,
		},
		{
			name:           "account with zero balance",
			accountID:      456,
			initialBalance: decimal.Zero,
			shouldError:    false,
		},
		{
			name:           "duplicate account ID should fail",
			accountID:      123,
			initialBalance: decimal.NewFromFloat(200.00),
			shouldError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Create(tc.accountID, tc.initialBalance)

			if tc.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify the account was created correctly
				account, err := repo.GetByID(tc.accountID)
				assert.NoError(t, err)
				assert.Equal(t, tc.accountID, account.ID)
				assert.True(t, tc.initialBalance.Equal(account.Balance))
				assert.False(t, account.CreatedAt.IsZero())
				assert.False(t, account.UpdatedAt.IsZero())
			}
		})
	}
}

func TestAccountRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Create test account
	accountID := int64(123)
	balance := decimal.NewFromFloat(100.50)
	err := repo.Create(accountID, balance)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		accountID   int64
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "existing account",
			accountID:   123,
			shouldError: false,
		},
		{
			name:        "non-existent account",
			accountID:   999,
			shouldError: true,
			errorMsg:    "account not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			account, err := repo.GetByID(tc.accountID)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, account)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, account)
				assert.Equal(t, tc.accountID, account.ID)
				assert.True(t, balance.Equal(account.Balance))
			}
		})
	}
}

func TestAccountRepository_UpdateBalance(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Create test account
	accountID := int64(123)
	initialBalance := decimal.NewFromFloat(100.50)
	err := repo.Create(accountID, initialBalance)
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
			name:        "non-existent account",
			accountID:   999,
			newBalance:  decimal.NewFromFloat(50.00),
			shouldError: true,
			errorMsg:    "account not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.UpdateBalance(tc.accountID, tc.newBalance)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)

				// Verify the balance was updated
				account, err := repo.GetByID(tc.accountID)
				assert.NoError(t, err)
				assert.True(t, tc.newBalance.Equal(account.Balance))
			}
		})
	}
}

func TestAccountRepository_Exists(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Create test account
	accountID := int64(123)
	balance := decimal.NewFromFloat(100.50)
	err := repo.Create(accountID, balance)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		accountID   int64
		shouldExist bool
	}{
		{
			name:        "existing account",
			accountID:   123,
			shouldExist: true,
		},
		{
			name:        "non-existent account",
			accountID:   999,
			shouldExist: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exists, err := repo.Exists(tc.accountID)
			assert.NoError(t, err)
			assert.Equal(t, tc.shouldExist, exists)
		})
	}
}
