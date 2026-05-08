package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type WalletTypeController struct {
	service *services.WalletTypeService
}

func NewWalletTypeController() *WalletTypeController {
	return &WalletTypeController{
		service: services.NewWalletTypeService(),
	}
}

func (c *WalletTypeController) CreateWalletType(ctx *gin.Context) {
	var req dtos.CreateWalletTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	walletType, err := c.service.CreateWalletType(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, walletType)
}

func (c *WalletTypeController) GetWalletTypes(ctx *gin.Context) {
	walletTypes, total, err := c.service.GetWalletTypes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": walletTypes, "total": total})
}

func (c *WalletTypeController) GetWalletType(ctx *gin.Context) {
	walletType, err := c.service.GetWalletType(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Wallet Type not found"})
		return
	}
	ctx.JSON(http.StatusOK, walletType)
}

func (c *WalletTypeController) UpdateWalletType(ctx *gin.Context) {
	var req dtos.UpdateWalletTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateWalletType(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Wallet Type updated successfully"})
}

func (c *WalletTypeController) DeleteWalletType(ctx *gin.Context) {
	if err := c.service.DeleteWalletType(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Wallet Type deleted successfully"})
}
