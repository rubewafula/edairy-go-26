package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
	"gorm.io/gorm"
)

type AccountCategoryController struct {
	service *services.AccountCategoryService
}

func NewAccountCategoryController() *AccountCategoryController {
	return &AccountCategoryController{service: services.NewAccountCategoryService()}
}

// Create godoc
// @Summary Create a new account category
// @Description Create a new account category with the provided details
// @Tags Account Categories
// @Accept json
// @Produce json
// @Param category body dtos.CreateAccountCategoryRequest true "Account category creation request"
// @Success 201 {object} models.AccountCategory
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /account-categories [post]
func (h *AccountCategoryController) Create(c *gin.Context) {
	var req dtos.CreateAccountCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint64)
	category, err := h.service.CreateAccountCategory(req, userID)
	if err != nil {
		log.Println("AccountCategoryController.Create Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account category"})
		return
	}
	c.JSON(http.StatusCreated, category)
}

// List godoc
// @Summary Get all account categories
// @Description Retrieve a list of all account categories with pagination
// @Tags Account Categories
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /account-categories [get]
func (h *AccountCategoryController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	categories, total, err := h.service.GetAccountCategories(page, limit)
	if err != nil {
		log.Println("AccountCategoryController.List Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve account categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories, "total": total})
}

// Get godoc
// @Summary Get an account category by ID
// @Description Retrieve a single account category by its ID
// @Tags Account Categories
// @Accept json
// @Produce json
// @Param id path string true "Account Category ID"
// @Success 200 {object} dtos.AccountCategoryResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /account-categories/{id} [get]
func (h *AccountCategoryController) Get(c *gin.Context) {
	id := c.Param("id")
	category, err := h.service.GetAccountCategory(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Account category with ID %s not found", id)})
			return
		}
		log.Println("AccountCategoryController.Get Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve account category"})
		return
	}
	c.JSON(http.StatusOK, category)
}

// Update godoc
// @Summary Update an existing account category
// @Description Update an account category with the provided details by ID
// @Tags Account Categories
// @Accept json
// @Produce json
// @Param id path string true "Account Category ID"
// @Param category body dtos.UpdateAccountCategoryRequest true "Account category update request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /account-categories/{id} [put]
func (h *AccountCategoryController) Update(c *gin.Context) {
	id := c.Param("id")
	var req dtos.UpdateAccountCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint64)
	err := h.service.UpdateAccountCategory(id, req, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Account category with ID %s not found", id)})
			return
		}
		log.Println("AccountCategoryController.Update Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account category updated successfully"})
}

// Delete godoc
// @Summary Delete an account category
// @Description Soft delete an account category by its ID
// @Tags Account Categories
// @Accept json
// @Produce json
// @Param id path string true "Account Category ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /account-categories/{id} [delete]
func (h *AccountCategoryController) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(uint64)
	err := h.service.DeleteAccountCategory(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Account category with ID %s not found", id)})
			return
		}
		log.Println("AccountCategoryController.Delete Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account category deleted successfully"})
}
