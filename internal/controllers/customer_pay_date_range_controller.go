package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"gorm.io/gorm"
)

type CustomerPayDateRangeController struct {
	service *services.CustomerPayDateRangeService
}

func NewCustomerPayDateRangeController() *CustomerPayDateRangeController {
	return &CustomerPayDateRangeController{
		service: services.NewCustomerPayDateRangeService(),
	}
}

// CreateCustomerType godoc
// @Summary Create a new customer type
// @Description Create a new customer type with the provided details
// @Tags Customer Types
// @Accept json
// @Produce json
// @Param customerType body dtos.CreateCustomerTypeRequest true "Customer Type object to be created"
// @Success 201 {object} models.CustomerType
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customer-types [post]
func (c *CustomerPayDateRangeController) CreateCustomerPayDateRange(ctx *gin.Context) {
	var req dtos.CreateCustomerPayDateRangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.MustGet("user_id").(uint64)

	customerPayDateRange, err := c.service.CreateCustomerPayDateRange(req, userID)
	if err != nil {
		log.Printf("[CustomerPayDateRangeController.CreateCustomerPayDateRange] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, customerPayDateRange)
}

// GetCustomerTypes godoc
// @Summary Get all customer types
// @Description Retrieve a list of all customer types
// @Tags Customer Types
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {array} dtos.CustomerTypeResponse
// @Failure 500 {object} map[string]string
// @Router /customer-types [get]
func (c *CustomerPayDateRangeController) GetCustomerPayDateRanges(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	customerTypes, total, err := c.service.GetCustomerPayDateRanges(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": customerTypes, "total": total})
}

// GetCustomerType godoc
// @Summary Get a single customer type by ID
// @Description Retrieve a customer type by its ID
// @Tags Customer Types
// @Produce json
// @Param id path string true "Customer Type ID"
// @Success 200 {object} dtos.CustomerTypeResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customer-types/{id} [get]
func (c *CustomerPayDateRangeController) GetCustomerPayDateRange(ctx *gin.Context) {
	id := ctx.Param("id")
	customerType, err := c.service.GetCustomerPayDateRange(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Customer type not found"})
			return
		}
		log.Printf("[CustomerPayDateRangeController.UpdateCustomerPayDateRange] Error updating range %s: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, customerType)
}

// UpdateCustomerType godoc
// @Summary Update an existing customer type
// @Description Update a customer type with the provided ID and details
// @Tags Customer Types
// @Accept json
// @Produce json
// @Param id path string true "Customer Type ID"
// @Param customerType body dtos.UpdateCustomerTypeRequest true "Updated customer type object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customer-types/{id} [put]
func (c *CustomerPayDateRangeController) UpdateCustomerPayDateRange(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateCustomerPayDateRangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.MustGet("user_id").(uint64)

	err := c.service.UpdateCustomerPayDateRange(id, req, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Customer type not found"})
			return
		}
		log.Printf("[CustomerPayDateRangeController.DeleteCustomerPayDateRange] Error deleting range %s: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Customer type updated successfully"})
}

// DeleteCustomerType godoc
// @Summary Delete a customer type
// @Description Delete a customer type by its ID
// @Tags Customer Types
// @Produce json
// @Param id path string true "Customer Type ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customer-types/{id} [delete]
func (c *CustomerPayDateRangeController) DeleteCustomerPayDateRange(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteCustomerPayDateRange(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Customer type not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Customer type deleted successfully"})
}
