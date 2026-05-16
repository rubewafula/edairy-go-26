package controllers

import (
	"net/http"

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
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	transfer, err := c.service.CreateShareTransfer(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, transfer)
}

func (c *ShareTransferController) GetShareTransfers(ctx *gin.Context) {
	transfers, total, err := c.service.GetShareTransfers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": transfers, "total": total})
}

func (c *ShareTransferController) GetShareTransfer(ctx *gin.Context) {
	transfer, err := c.service.GetShareTransfer(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Share transfer not found"})
		return
	}
	ctx.JSON(http.StatusOK, transfer)
}

func (c *ShareTransferController) UpdateShareTransfer(ctx *gin.Context) {
	var req dtos.UpdateShareTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateShareTransfer(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share transfer updated successfully"})
}

func (c *ShareTransferController) DeleteShareTransfer(ctx *gin.Context) {
	if err := c.service.DeleteShareTransfer(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share transfer deleted successfully"})
}
