package main

import (
	"log"
	"os"

	"internal-transfer-system/internal/database"
)

func main() {
	// Load database configuration
	config := database.NewConfig()

	// Connect to database
	if err := database.Connect(config); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create database tables
	if err := database.CreateTables(database.DB); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	log.Println("Internal Transfer System initialized successfully")
	log.Println("Database connected and tables created")

	// For now, just keep the application running
	// In next iteration, we'll add the HTTP server
	select {}
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
