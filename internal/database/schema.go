package database

import (
	"fmt"
	"log"

	"internal-transfer-system/internal/model"

	"gorm.io/gorm"
)

// CreateTables creates the necessary tables using GORM auto-migration
func CreateTables(db *gorm.DB) error {
	// Auto-migrate the models
	err := db.AutoMigrate(
		&model.Account{},
		&model.Transaction{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate tables: %w", err)
	}

	log.Println("Database tables created successfully")
	return nil
}

// DropTables drops all tables (useful for testing)
func DropTables(db *gorm.DB) error {
	// Drop tables in reverse order to respect foreign key constraints
	err := db.Migrator().DropTable(&model.Transaction{}, &model.Account{})
	if err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	log.Println("Database tables dropped successfully")
	return nil
}
