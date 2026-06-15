package controllers

import (
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type ShareTypeController struct {
	service *services.ShareTypeService
}

func NewShareTypeController() *ShareTypeController {
	return &ShareTypeController{
		service: services.NewShareTypeService(),
	}
}

func (c *ShareTypeController) CreateShareType(ctx *gin.Context) {
	var req dtos.CreateShareTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[ShareTypeController.CreateShareType] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[ShareTypeController.CreateShareType] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	shareType, err := c.service.CreateShareType(req)
	if err != nil {
		log.Printf("[ShareTypeController.CreateShareType] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create share type"})
		return
	}
	ctx.JSON(http.StatusCreated, shareType)
}

func (c *ShareTypeController) GetShareTypes(ctx *gin.Context) {
	shareTypes, total, err := c.service.GetShareTypes()
	if err != nil {
		log.Printf("[ShareTypeController.GetShareTypes] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve share types"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": shareTypes, "total": total})
}

func (c *ShareTypeController) GetShareType(ctx *gin.Context) {
	shareType, err := c.service.GetShareType(ctx.Param("id"))
	if err != nil {
		log.Printf("[ShareTypeController.GetShareType] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Share type not found"})
		return
	}
	ctx.JSON(http.StatusOK, shareType)
}

func (c *ShareTypeController) UpdateShareType(ctx *gin.Context) {
	var req dtos.UpdateShareTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[ShareTypeController.UpdateShareType] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[ShareTypeController.UpdateShareType] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateShareType(ctx.Param("id"), req); err != nil {
		log.Printf("[ShareTypeController.UpdateShareType] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update share type"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share type updated successfully"})
}

func (c *ShareTypeController) DeleteShareType(ctx *gin.Context) {
	if err := c.service.DeleteShareType(ctx.Param("id")); err != nil {
		log.Printf("[ShareTypeController.DeleteShareType] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete share type"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share type deleted successfully"})
}
