package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"internal-transfer-system/internal/database"
	"internal-transfer-system/internal/router"
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

	// Setup HTTP router
	r := router.SetupRouter(database.DB)

	// Get server port from environment or use default
	port := getEnv("PORT", "8080")

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Internal Transfer System started successfully")
	log.Printf("Server running on http://localhost:%s", port)
	log.Println("API endpoints:")
	log.Println("  POST /accounts - Create account")
	log.Println("  GET /accounts/{account_id} - Get account balance")
	log.Println("  POST /transactions - Create transaction")
	log.Println("  GET /health - Health check")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
