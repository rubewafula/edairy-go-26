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

type AccountSubAccountController struct {
	service *services.AccountSubAccountService
}

func NewAccountSubAccountController() *AccountSubAccountController {
	return &AccountSubAccountController{service: services.NewAccountSubAccountService()}
}

// CreateAccountSubAccount godoc
// @Summary Create a new account sub-account
// @Description Create a new account sub-account with the provided details
// @Tags Account Sub-Accounts
// @Accept json
// @Produce json
// @Param sub_account body dtos.CreateAccountSubAccountRequest true "Account sub-account creation request"
// @Success 201 {object} models.AccountSubAccount
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /account-sub-accounts [post]
func (h *AccountSubAccountController) CreateAccountSubAccount(c *gin.Context) {
	var req dtos.CreateAccountSubAccountRequest
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
	subAccount, err := h.service.CreateAccountSubAccount(req, userID)
	if err != nil { // Log the exact error but return a user-friendly message
		log.Println("AccountSubAccountController.CreateAccountSubAccount Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, subAccount)
}

// GetAccountSubAccounts godoc
// @Summary Get all account sub-accounts
// @Description Retrieve a list of all account sub-accounts with pagination and optional filtering by account ID
// @Tags Account Sub-Accounts
// @Accept json
// @Produce json
// @Param account_id query string false "Filter by Account ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} dtos.PaginatedResponse{data=[]dtos.AccountSubAccountResponse}
// @Failure 500 {object} map[string]string
// @Router /account-sub-accounts [get]
func (h *AccountSubAccountController) GetAccountSubAccounts(c *gin.Context) {
	accountID := c.Query("account_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	subAccounts, total, err := h.service.GetAccountSubAccounts(accountID, page, limit)
	if err != nil { // Log the exact error but return a user-friendly message
		log.Println("AccountSubAccountController.GetAccountSubAccounts Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": subAccounts, "total": total})
}

// GetAccountSubAccount godoc
// @Summary Get an account sub-account by ID
// @Description Retrieve a single account sub-account by its ID
// @Tags Account Sub-Accounts
// @Accept json
// @Produce json
// @Param id path string true "Account Sub-Account ID"
// @Success 200 {object} dtos.AccountSubAccountResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /account-sub-accounts/{id} [get]
func (h *AccountSubAccountController) GetAccountSubAccount(c *gin.Context) {
	id := c.Param("id")
	subAccount, err := h.service.GetAccountSubAccount(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Account Sub-Account with ID %s not found", id)})
			return
		}
		log.Println("AccountSubAccountController.GetAccountSubAccount Error:", err)                      // Log the exact error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve account sub-account"}) // User-friendly message
		return
	}
	c.JSON(http.StatusOK, subAccount)
}

// UpdateAccountSubAccount godoc
// @Summary Update an existing account sub-account
// @Description Update an account sub-account with the provided details by ID
// @Tags Account Sub-Accounts
// @Accept json
// @Produce json
// @Param id path string true "Account Sub-Account ID"
// @Param sub_account body dtos.UpdateAccountSubAccountRequest true "Account sub-account update request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /account-sub-accounts/{id} [put]
func (h *AccountSubAccountController) UpdateAccountSubAccount(c *gin.Context) {
	id := c.Param("id")
	var req dtos.UpdateAccountSubAccountRequest
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
	err := h.service.UpdateAccountSubAccount(id, req, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Account Sub-Account with ID %s not found", id)})
			return
		}
		log.Println("AccountSubAccountController.UpdateAccountSubAccount Error:", err)                 // Log the exact error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account sub-account"}) // User-friendly message
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account Sub-Account updated successfully"})
}

// DeleteAccountSubAccount godoc
// @Summary Delete an account sub-account
// @Description Soft delete an account sub-account by its ID
// @Tags Account Sub-Accounts
// @Accept json
// @Produce json
// @Param id path string true "Account Sub-Account ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /account-sub-accounts/{id} [delete]
func (h *AccountSubAccountController) DeleteAccountSubAccount(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(uint64)
	err := h.service.DeleteAccountSubAccount(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Account Sub-Account with ID %s not found", id)})
			return
		}
		log.Println("AccountSubAccountController.DeleteAccountSubAccount Error:", err)                 // Log the exact error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account sub-account"}) // User-friendly message
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account Sub-Account deleted successfully"})
}
