package handlers

import (
	"billing-service/internal/dto"
	"billing-service/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	billingService *services.BillingService
}

func NewHandlers(billingService *services.BillingService) *Handlers {
	return &Handlers{
		billingService: billingService,
	}
}

func (h *Handlers) CreateAccount(c *gin.Context) {
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	account, err := h.billingService.CreateAccount(req.UserID)
	if err != nil {
		if err == services.ErrAccountExists {
			c.JSON(409, gin.H{"error": "Account already exists"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to create account"})
		return
	}

	response := dto.CreateAccountResponse{
		AccountID: account.ID,
		Balance:   account.Balance,
	}

	c.JSON(201, response)
}

func (h *Handlers) WithdrawMoney(c *gin.Context) {
	var req dto.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := h.billingService.WithdrawMoney(req.UserID, req.Amount)
	if err != nil {
		if err == services.ErrAccountNotFound {
			c.JSON(404, gin.H{"error": "Account not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to process withdrawal"})
		return
	}

	response := dto.WithdrawResponse{
		Success:    result.Success,
		NewBalance: result.NewBalance,
		Message:    result.Message,
	}

	statusCode := 200
	if !result.Success {
		statusCode = 400
	}

	c.JSON(statusCode, response)
}

func (h *Handlers) DepositMoney(c *gin.Context) {
	var req dto.DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := h.billingService.DepositMoney(req.UserID, req.Amount)
	if err != nil {
		if err == services.ErrAccountNotFound {
			c.JSON(404, gin.H{"error": "Account not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to process deposit"})
		return
	}

	response := dto.DepositResponse{
		Success:    result.Success,
		NewBalance: result.NewBalance,
		Message:    result.Message,
	}

	c.JSON(200, response)
}

func (h *Handlers) GetBalance(c *gin.Context) {
	userID := c.Param("user_id")

	account, err := h.billingService.GetAccount(userID)
	if err != nil {
		if err == services.ErrAccountNotFound {
			c.JSON(404, gin.H{"error": "Account not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to get account"})
		return
	}

	response := dto.BalanceResponse{
		UserID:  account.UserID,
		Balance: account.Balance,
	}
	c.JSON(200, response)
}

func (h *Handlers) GetPayments(c *gin.Context) {
	userID := c.Param("user_id")
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	payments, err := h.billingService.GetPayments(userID, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch payments"})
		return
	}

	c.JSON(200, gin.H{"payments": payments})
}

func (h *Handlers) Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
