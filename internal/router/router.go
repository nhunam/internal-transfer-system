package router

import (
	"database/sql"

	"internal-transfer-system/internal/handler"
	"internal-transfer-system/internal/repository"
	"internal-transfer-system/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the HTTP routes and returns a Gin router
func SetupRouter(db *sql.DB) *gin.Engine {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize repositories
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Initialize services
	accountService := service.NewAccountService(accountRepo)
	transactionService := service.NewTransactionService(db, transactionRepo, accountService)

	// Initialize handlers
	accountHandler := handler.NewAccountHandler(accountService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// Setup routes
	// Account routes
	router.POST("/accounts", accountHandler.CreateAccount)
	router.GET("/accounts/:account_id", accountHandler.GetAccount)

	// Transaction routes
	router.POST("/transactions", transactionHandler.CreateTransaction)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "internal-transfer-system",
		})
	})

	return router
}
