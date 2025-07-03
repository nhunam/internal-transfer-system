package handler

import (
	"net/http"

	"internal-transfer-system/internal/model"
	"internal-transfer-system/internal/service"
	"internal-transfer-system/internal/utils"

	"github.com/gin-gonic/gin"
)

// TransactionHandler handles HTTP requests for transaction operations
type TransactionHandler struct {
	transactionService *service.TransactionService
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// CreateTransaction handles POST /transactions
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var request model.CreateTransactionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if err := h.transactionService.CreateTransaction(&request); err != nil {
		// Determine appropriate HTTP status code based on error type
		statusCode := http.StatusInternalServerError

		// Check for specific business logic errors
		errorMessage := err.Error()
		if utils.ContainsAny(errorMessage, []string{
			"source account ID must be positive",
			"destination account ID must be positive",
			"source and destination accounts cannot be the same",
			"amount is required",
			"invalid amount format",
			"amount must be positive",
			"account validation failed",
			"account not found",
			"insufficient balance",
		}) {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error": errorMessage,
		})
		return
	}

	// Return empty response on success
	c.Status(http.StatusCreated)
}
