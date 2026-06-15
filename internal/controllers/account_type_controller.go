package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"gorm.io/gorm"
)

type AccountTypeController struct {
	service *services.AccountTypeService
}

func NewAccountTypeController() *AccountTypeController {
	return &AccountTypeController{service: services.NewAccountTypeService()}
}

// Create godoc
// @Summary Create a new account type
// @Description Create a new account type with the provided details
// @Tags Account Types
// @Accept json
// @Produce json
// @Param account_type body dtos.CreateAccountTypeRequest true "Account type creation request"
// @Success 201 {object} models.AccountType
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /account-types [post]
func (h *AccountTypeController) Create(c *gin.Context) {
	var req dtos.CreateAccountTypeRequest
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
	accountType, err := h.service.CreateAccountType(req, userID)
	if err != nil {
		log.Println("AccountTypeController.Create Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account type"})
		return
	}
	c.JSON(http.StatusCreated, accountType)
}

// List godoc
// @Summary Get all account types
// @Description Retrieve a list of all account types with pagination
// @Tags Account Types
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /account-types [get]
func (h *AccountTypeController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	accountTypes, total, err := h.service.GetAccountTypes(page, limit)
	if err != nil {
		log.Println("AccountTypeController.List Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve account types"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": accountTypes, "total": total})
}

// Get godoc
// @Summary Get an account type by ID
// @Description Retrieve a single account type by its ID
// @Tags Account Types
// @Accept json
// @Produce json
// @Param id path string true "Account Type ID"
// @Success 200 {object} dtos.AccountTypeResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /account-types/{id} [get]
func (h *AccountTypeController) Get(c *gin.Context) {
	id := c.Param("id")
	accountType, err := h.service.GetAccountType(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Account type with ID %s not found", id)})
			return
		}
		log.Println("AccountTypeController.Get Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve account type"})
		return
	}
	c.JSON(http.StatusOK, accountType)
}

// Update godoc
// @Summary Update an existing account type
// @Description Update an account type with the provided details by ID
// @Tags Account Types
// @Accept json
// @Produce json
// @Param id path string true "Account Type ID"
// @Param account_type body dtos.UpdateAccountTypeRequest true "Account type update request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /account-types/{id} [put]
func (h *AccountTypeController) Update(c *gin.Context) {
	id := c.Param("id")
	var req dtos.UpdateAccountTypeRequest
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
	err := h.service.UpdateAccountType(id, req, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Account type with ID %s not found", id)})
			return
		}
		log.Println("AccountTypeController.Update Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account type"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account type updated successfully"})
}

// Delete godoc
// @Summary Delete an account type
// @Description Soft delete an account type by its ID
// @Tags Account Types
// @Accept json
// @Produce json
// @Param id path string true "Account Type ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /account-types/{id} [delete]
func (h *AccountTypeController) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(uint64)
	err := h.service.DeleteAccountType(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Account type with ID %s not found", id)})
			return
		}
		log.Println("AccountTypeController.Delete Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account type"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account type deleted successfully"})
}
