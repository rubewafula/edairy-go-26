package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type CanMovementController struct {
	service *services.CanMovementService
}

func NewCanMovementController() *CanMovementController {
	return &CanMovementController{
		service: services.NewCanMovementService(),
	}
}

func (c *CanMovementController) CreateMovement(ctx *gin.Context) {
	var req dtos.CreateCanMovementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	movement, err := c.service.CreateMovement(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, _ := c.service.GetMovement(utils.Uint64ToString(movement.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *CanMovementController) GetMovements(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	movements, total, err := c.service.GetMovements(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": movements, "total": total})
}

func (c *CanMovementController) GetMovement(ctx *gin.Context) {
	movement, err := c.service.GetMovement(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Can movement not found"})
		return
	}
	ctx.JSON(http.StatusOK, movement)
}

func (c *CanMovementController) UpdateMovement(ctx *gin.Context) {
	var req dtos.UpdateCanMovementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateMovement(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Can movement updated successfully"})
}

func (c *CanMovementController) DeleteMovement(ctx *gin.Context) {
	if err := c.service.DeleteMovement(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Can movement deleted successfully"})
}
