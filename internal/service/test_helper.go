package service

import (
	"internal-transfer-system/internal/model"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		// Skip default transaction for better performance and concurrency
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

	// Configure SQLite for better concurrency
	sqlDB, err := db.DB()
	require.NoError(t, err)

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)

	// Auto-migrate the schema
	err = db.AutoMigrate(&model.Account{}, &model.Transaction{})
	require.NoError(t, err)

	return db
}
