package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"gorm.io/gorm"
)

type AccountController struct {
	service *services.AccountService
}

func NewAccountController() *AccountController {
	return &AccountController{service: services.NewAccountService()}
}

// CreateAccount godoc
// @Summary Create a new account
// @Description Create a new account with the provided details
// @Tags Accounts
// @Accept json
// @Produce json
// @Param account body dtos.CreateAccountRequest true "Account creation request"
// @Success 201 {object} models.Account
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts [post]
func (h *AccountController) CreateAccount(c *gin.Context) {
	var req dtos.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint64)
	account, err := h.service.CreateAccount(req, userID)
	if err != nil { // Log the exact error but return a user-friendly message
		log.Println("AccountController.CreateAccount Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, account)
}

// GetAccounts godoc
// @Summary Get all accounts
// @Description Retrieve a list of all accounts with pagination
// @Tags Accounts
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} dtos.PaginatedResponse{data=[]dtos.AccountResponse}
// @Failure 500 {object} map[string]string
// @Router /accounts [get]
func (h *AccountController) GetAccounts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	accounts, total, err := h.service.GetAccounts(page, limit)
	if err != nil { // Log the exact error but return a user-friendly message
		log.Println("AccountController.GetAccounts Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": accounts, "total": total})
}

// GetAccount godoc
// @Summary Get an account by ID
// @Description Retrieve a single account by its ID
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} dtos.AccountResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts/{id} [get]
func (h *AccountController) GetAccount(c *gin.Context) {
	id := c.Param("id")
	account, err := h.service.GetAccount(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Account with ID %s not found", id)})
			return
		}
		log.Println("AccountController.GetAccount Error:", err)                              // Log the exact error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve account"}) // User-friendly message
		return
	}
	c.JSON(http.StatusOK, account)
}

// UpdateAccount godoc
// @Summary Update an existing account
// @Description Update an account with the provided details by ID
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param account body dtos.UpdateAccountRequest true "Account update request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts/{id} [put]
func (h *AccountController) UpdateAccount(c *gin.Context) {
	id := c.Param("id")
	var req dtos.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint64)
	err := h.service.UpdateAccount(id, req, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Account with ID %s not found", id)})
			return
		}
		log.Println("AccountController.UpdateAccount Error:", err)                         // Log the exact error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account"}) // User-friendly message
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account updated successfully"})
}

// DeleteAccount godoc
// @Summary Delete an account
// @Description Soft delete an account by its ID
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts/{id} [delete]
func (h *AccountController) DeleteAccount(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(uint64)
	err := h.service.DeleteAccount(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Account with ID %s not found", id)})
			return
		}
		log.Println("AccountController.DeleteAccount Error:", err)                         // Log the exact error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"}) // User-friendly message
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

// GetTrialBalance godoc
// @Summary Get the Trial Balance report
// @Description Retrieve the current balances for all accounts with activity in the general ledger
// @Tags Accounting Reports
// @Produce json
// @Success 200 {object} dtos.TrialBalanceResponse
// @Failure 500 {object} map[string]string
// @Router /accounting/trial-balance [get]
func (h *AccountController) GetTrialBalance(c *gin.Context) {
	result, err := h.service.GetTrialBalance()
	if err != nil {
		// Log exact error but return user-friendly message
		log.Println("AccountController.GetTrialBalance Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate trial balance"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetProfitLoss godoc
// @Summary Get the Profit and Loss report
// @Description Retrieve the revenues and expenses summarized from the general ledger
// @Tags Accounting Reports
// @Produce json
// @Success 200 {object} dtos.ProfitLossResponse
// @Failure 500 {object} map[string]string
// @Router /accounting/profit-loss [get]
func (h *AccountController) GetProfitLoss(c *gin.Context) {
	result, err := h.service.GetProfitLoss()
	if err != nil {
		// Log exact error but return user-friendly message
		log.Println("AccountController.GetProfitLoss Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate profit and loss report"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetBalanceSheet godoc
// @Summary Get the Balance Sheet report
// @Description Retrieve the current position of Assets, Liabilities and Equity
// @Tags Accounting Reports
// @Produce json
// @Success 200 {object} dtos.BalanceSheetResponse
// @Failure 500 {object} map[string]string
// @Router /accounting/balance-sheet [get]
func (h *AccountController) GetBalanceSheet(c *gin.Context) {
	result, err := h.service.GetBalanceSheet()
	if err != nil {
		// Log exact error but return user-friendly message
		log.Println("AccountController.GetBalanceSheet Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate balance sheet report"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}
