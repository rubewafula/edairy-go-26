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

type ShareTransferController struct {
	service *services.ShareTransferService
}

func NewShareTransferController() *ShareTransferController {
	return &ShareTransferController{
		service: services.NewShareTransferService(),
	}
}

func (c *ShareTransferController) CreateShareTransfer(ctx *gin.Context) {
	var req dtos.CreateShareTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[ShareTransferController.CreateShareTransfer] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[ShareTransferController.CreateShareTransfer] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	transfer, err := c.service.CreateShareTransfer(req)
	if err != nil {
		log.Printf("[ShareTransferController.CreateShareTransfer] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create share transfer"})
		return
	}
	ctx.JSON(http.StatusCreated, transfer)
}

func (c *ShareTransferController) GetShareTransfers(ctx *gin.Context) {
	transfers, total, err := c.service.GetShareTransfers()
	if err != nil {
		log.Printf("[ShareTransferController.GetShareTransfers] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve share transfers"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": transfers, "total": total})
}

func (c *ShareTransferController) GetShareTransfer(ctx *gin.Context) {
	transfer, err := c.service.GetShareTransfer(ctx.Param("id"))
	if err != nil {
		log.Printf("[ShareTransferController.GetShareTransfer] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Share transfer not found"})
		return
	}
	ctx.JSON(http.StatusOK, transfer)
}

func (c *ShareTransferController) UpdateShareTransfer(ctx *gin.Context) {
	var req dtos.UpdateShareTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[ShareTransferController.UpdateShareTransfer] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[ShareTransferController.UpdateShareTransfer] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateShareTransfer(ctx.Param("id"), req); err != nil {
		log.Printf("[ShareTransferController.UpdateShareTransfer] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update share transfer"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share transfer updated successfully"})
}

func (c *ShareTransferController) DeleteShareTransfer(ctx *gin.Context) {
	if err := c.service.DeleteShareTransfer(ctx.Param("id")); err != nil {
		log.Printf("[ShareTransferController.DeleteShareTransfer] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete share transfer"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share transfer deleted successfully"})
}
