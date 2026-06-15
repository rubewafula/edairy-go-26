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

type TransactionPostingRuleController struct {
	service *services.TransactionPostingRuleService
}

func NewTransactionPostingRuleController() *TransactionPostingRuleController {
	return &TransactionPostingRuleController{service: services.NewTransactionPostingRuleService()}
}

// CreateTransactionPostingRule godoc
// @Summary Create a new transaction posting rule
// @Description Create a new transaction posting rule with the provided details
// @Tags Transaction Posting Rules
// @Accept json
// @Produce json
// @Param rule body dtos.CreateTransactionPostingRuleRequest true "Transaction posting rule creation request"
// @Success 201 {object} models.TransactionPostingRule
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transaction-posting-rules [post]
func (h *TransactionPostingRuleController) CreateTransactionPostingRule(c *gin.Context) {
	var req dtos.CreateTransactionPostingRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("TransactionPostingRuleController.CreateTransactionPostingRule Binding Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		log.Println("TransactionPostingRuleController.CreateTransactionPostingRule Validation Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint64)
	rule, err := h.service.CreateTransactionPostingRule(req, userID)
	if err != nil {
		log.Println("TransactionPostingRuleController.CreateTransactionPostingRule Service Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction posting rule"})
		return
	}
	c.JSON(http.StatusCreated, rule)
}

// GetTransactionPostingRules godoc
// @Summary Get all transaction posting rules
// @Description Retrieve a list of all transaction posting rules with pagination
// @Tags Transaction Posting Rules
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} dtos.PaginatedResponse{data=[]dtos.TransactionPostingRuleResponse}
// @Failure 500 {object} map[string]string
// @Router /transaction-posting-rules [get]
func (h *TransactionPostingRuleController) GetTransactionPostingRules(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	rules, total, err := h.service.GetTransactionPostingRules(page, limit)
	if err != nil {
		log.Println("TransactionPostingRuleController.GetTransactionPostingRules Service Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transaction posting rules"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rules, "total": total})
}

// GetTransactionPostingRule godoc
// @Summary Get a transaction posting rule by ID
// @Description Retrieve a single transaction posting rule by its ID
// @Tags Transaction Posting Rules
// @Accept json
// @Produce json
// @Param id path string true "Transaction Posting Rule ID"
// @Success 200 {object} dtos.TransactionPostingRuleResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transaction-posting-rules/{id} [get]
func (h *TransactionPostingRuleController) GetTransactionPostingRule(c *gin.Context) {
	id := c.Param("id")
	rule, err := h.service.GetTransactionPostingRule(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Transaction posting rule with ID %s not found", id)})
			return
		}
		log.Println("TransactionPostingRuleController.GetTransactionPostingRule Service Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transaction posting rule"})
		return
	}
	c.JSON(http.StatusOK, rule)
}

// UpdateTransactionPostingRule godoc
// @Summary Update an existing transaction posting rule
// @Description Update a transaction posting rule with the provided details by ID
// @Tags Transaction Posting Rules
// @Accept json
// @Produce json
// @Param id path string true "Transaction Posting Rule ID"
// @Param rule body dtos.UpdateTransactionPostingRuleRequest true "Transaction posting rule update request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transaction-posting-rules/{id} [put]
func (h *TransactionPostingRuleController) UpdateTransactionPostingRule(c *gin.Context) {
	id := c.Param("id")
	var req dtos.UpdateTransactionPostingRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("TransactionPostingRuleController.UpdateTransactionPostingRule Binding Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		log.Println("TransactionPostingRuleController.UpdateTransactionPostingRule Validation Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint64)
	err := h.service.UpdateTransactionPostingRule(id, req, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Transaction posting rule with ID %s not found", id)})
			return
		}
		log.Println("TransactionPostingRuleController.UpdateTransactionPostingRule Service Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction posting rule"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaction posting rule updated successfully"})
}

// DeleteTransactionPostingRule godoc
// @Summary Delete a transaction posting rule
// @Description Soft delete a transaction posting rule by its ID
// @Tags Transaction Posting Rules
// @Accept json
// @Produce json
// @Param id path string true "Transaction Posting Rule ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transaction-posting-rules/{id} [delete]
func (h *TransactionPostingRuleController) DeleteTransactionPostingRule(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(uint64)
	err := h.service.DeleteTransactionPostingRule(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Transaction posting rule with ID %s not found", id)})
			return
		}
		log.Println("TransactionPostingRuleController.DeleteTransactionPostingRule Service Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction posting rule"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaction posting rule deleted successfully"})
}
