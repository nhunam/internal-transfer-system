package database

import (
	"database/sql"
	"fmt"
	"log"
)

// CreateTables creates the necessary tables for the application
func CreateTables(db *sql.DB) error {
	// Create accounts table
	accountsTable := `
		CREATE TABLE IF NOT EXISTS accounts (
			account_id BIGINT PRIMARY KEY,
			balance DECIMAL(20, 8) NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	// Create transactions table
	transactionsTable := `
		CREATE TABLE IF NOT EXISTS transactions (
			transaction_id BIGSERIAL PRIMARY KEY,
			source_account_id BIGINT NOT NULL,
			destination_account_id BIGINT NOT NULL,
			amount DECIMAL(20, 8) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (source_account_id) REFERENCES accounts(account_id),
			FOREIGN KEY (destination_account_id) REFERENCES accounts(account_id)
		);
	`

	// Create indexes for better performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_transactions_source_account ON transactions(source_account_id);",
		"CREATE INDEX IF NOT EXISTS idx_transactions_destination_account ON transactions(destination_account_id);",
		"CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);",
		"CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);",
	}

	// Create trigger function for updating updated_at timestamps
	updateTriggerFunction := `
		CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
		$$ language 'plpgsql';
	`

	// Create triggers for updating updated_at
	triggers := []string{
		`CREATE TRIGGER IF NOT EXISTS update_accounts_updated_at 
		 BEFORE UPDATE ON accounts 
		 FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();`,
		`CREATE TRIGGER IF NOT EXISTS update_transactions_updated_at 
		 BEFORE UPDATE ON transactions 
		 FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();`,
	}

	// Execute all SQL statements
	statements := []string{accountsTable, transactionsTable}
	statements = append(statements, indexes...)
	statements = append(statements, updateTriggerFunction)
	statements = append(statements, triggers...)

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	log.Println("Database tables created successfully")
	return nil
}

// DropTables drops all tables (useful for testing)
func DropTables(db *sql.DB) error {
	dropStatements := []string{
		"DROP TABLE IF EXISTS transactions CASCADE;",
		"DROP TABLE IF EXISTS accounts CASCADE;",
		"DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;",
	}

	for _, stmt := range dropStatements {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to drop table: %w", err)
		}
	}

	log.Println("Database tables dropped successfully")
	return nil
}
