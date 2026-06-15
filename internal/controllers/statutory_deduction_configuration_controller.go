package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
	"gorm.io/gorm"
)

type StatutoryDeductionConfigurationController struct {
	service *services.StatutoryDeductionConfigurationService
}

func NewStatutoryDeductionConfigurationController() *StatutoryDeductionConfigurationController {
	return &StatutoryDeductionConfigurationController{
		service: services.NewStatutoryDeductionConfigurationService(),
	}
}

// CreateStatutoryDeductionConfiguration godoc
// @Summary Create a new statutory deduction configuration
// @Description Create a new statutory deduction configuration with the provided details
// @Tags Statutory Deductions
// @Accept json
// @Produce json
// @Param config body dtos.CreateStatutoryDeductionConfigurationRequest true "Statutory Deduction Configuration object to be created"
// @Success 201 {object} models.StatutoryDeductionConfiguration
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /statutory-deduction-configurations [post]
func (c *StatutoryDeductionConfigurationController) CreateConfiguration(ctx *gin.Context) {
	var req dtos.CreateStatutoryDeductionConfigurationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.MustGet("user_id").(uint64)

	config, err := c.service.CreateConfiguration(req, userID)
	if err != nil {
		log.Printf("[StatutoryDeductionConfigurationController.CreateConfiguration] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, config)
}

// GetStatutoryDeductionConfigurations godoc
// @Summary Get all statutory deduction configurations
// @Description Retrieve a list of all statutory deduction configurations
// @Tags Statutory Deductions
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {array} dtos.StatutoryDeductionConfigurationResponse
// @Failure 500 {object} map[string]string
// @Router /statutory-deduction-configurations [get]
func (c *StatutoryDeductionConfigurationController) GetConfigurations(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	configs, total, err := c.service.GetConfigurations(page, limit)
	if err != nil {
		log.Printf("[StatutoryDeductionConfigurationController.GetConfigurations] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve statutory deduction configurations"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": configs, "total": total})
}

// GetStatutoryDeductionConfiguration godoc
// @Summary Get a single statutory deduction configuration by ID
// @Description Retrieve a statutory deduction configuration by its ID
// @Tags Statutory Deductions
// @Produce json
// @Param id path string true "Statutory Deduction Configuration ID"
// @Success 200 {object} dtos.StatutoryDeductionConfigurationResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /statutory-deduction-configurations/{id} [get]
func (c *StatutoryDeductionConfigurationController) GetConfiguration(ctx *gin.Context) {
	id := ctx.Param("id")
	config, err := c.service.GetConfiguration(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Statutory deduction configuration not found"})
			return
		}
		log.Printf("[StatutoryDeductionConfigurationController.GetConfiguration] Error getting config %s: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, config)
}

// UpdateStatutoryDeductionConfiguration godoc
// @Summary Update an existing statutory deduction configuration
// @Description Update a statutory deduction configuration with the provided ID and details
// @Tags Statutory Deductions
// @Accept json
// @Produce json
// @Param id path string true "Statutory Deduction Configuration ID"
// @Param config body dtos.UpdateStatutoryDeductionConfigurationRequest true "Updated statutory deduction configuration object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /statutory-deduction-configurations/{id} [put]
func (c *StatutoryDeductionConfigurationController) UpdateConfiguration(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateStatutoryDeductionConfigurationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[StatutoryDeductionConfigurationController.UpdateConfiguration] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[StatutoryDeductionConfigurationController.UpdateConfiguration] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.MustGet("user_id").(uint64)

	err := c.service.UpdateConfiguration(id, req, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Statutory deduction configuration not found"})
			return
		}
		log.Printf("[StatutoryDeductionConfigurationController.UpdateConfiguration] Error updating config %s: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Statutory deduction configuration updated successfully"})
}

// DeleteStatutoryDeductionConfiguration godoc
// @Summary Delete a statutory deduction configuration
// @Description Delete a statutory deduction configuration by its ID
// @Tags Statutory Deductions
// @Produce json
// @Param id path string true "Statutory Deduction Configuration ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /statutory-deduction-configurations/{id} [delete]
func (c *StatutoryDeductionConfigurationController) DeleteConfiguration(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteConfiguration(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Statutory deduction configuration not found"})
			return
		}
		log.Printf("[StatutoryDeductionConfigurationController.DeleteConfiguration] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete statutory deduction configuration"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Statutory deduction configuration deleted successfully"})
}
