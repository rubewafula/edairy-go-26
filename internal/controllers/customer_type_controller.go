package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
	"gorm.io/gorm"
)

type CustomerTypeController struct {
	service *services.CustomerTypeService
	// Removed validator here as it's not used directly in the controller methods
}

func NewCustomerTypeController() *CustomerTypeController {
	return &CustomerTypeController{
		service: services.NewCustomerTypeService(),
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
func (c *CustomerTypeController) CreateCustomerType(ctx *gin.Context) {
	var req dtos.CreateCustomerTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Added validation as per other controllers
	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}
	userID := ctx.MustGet("user_id").(uint64)

	customerType, err := c.service.CreateCustomerType(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, customerType)
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
func (c *CustomerTypeController) GetCustomerTypes(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))
	customerTypes, total, err := c.service.GetCustomerTypes(page, limit)
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
func (c *CustomerTypeController) GetCustomerType(ctx *gin.Context) {
	id := ctx.Param("id")
	customerType, err := c.service.GetCustomerType(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Customer type not found"})
			return
		}
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
func (c *CustomerTypeController) UpdateCustomerType(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateCustomerTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Added validation as per other controllers
	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}
	userID := ctx.MustGet("user_id").(uint64)

	err := c.service.UpdateCustomerType(id, req, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Customer type not found"})
			return
		}
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
func (c *CustomerTypeController) DeleteCustomerType(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteCustomerType(id)
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
