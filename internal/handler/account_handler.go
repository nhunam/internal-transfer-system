package handler

import (
	"net/http"
	"strconv"

	"internal-transfer-system/internal/model"
	"internal-transfer-system/internal/service"
	"internal-transfer-system/internal/utils"

	"github.com/gin-gonic/gin"
)

// AccountHandler handles HTTP requests for account operations
type AccountHandler struct {
	accountService *service.AccountService
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// CreateAccount handles POST /accounts
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var request model.CreateAccountRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if err := h.accountService.CreateAccount(&request); err != nil {
		// Determine appropriate HTTP status code based on error type
		statusCode := http.StatusInternalServerError

		// Check for specific business logic errors
		errorMessage := err.Error()
		if utils.ContainsAny(errorMessage, []string{"account already exists", "account ID must be positive", "initial balance cannot be negative", "invalid initial balance format"}) {
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

// GetAccount handles GET /accounts/{account_id}
func (h *AccountHandler) GetAccount(c *gin.Context) {
	accountIDStr := c.Param("account_id")
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid account ID format",
		})
		return
	}

	account, err := h.accountService.GetAccount(accountID)
	if err != nil {
		// Determine appropriate HTTP status code based on error type
		statusCode := http.StatusInternalServerError

		// Check for specific business logic errors
		errorMessage := err.Error()
		if utils.ContainsAny(errorMessage, []string{"account not found", "account ID must be positive"}) {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"error": errorMessage,
		})
		return
	}

	c.JSON(http.StatusOK, account)
}
