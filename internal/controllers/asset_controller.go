package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type AssetController struct {
	service *services.AssetService
}

func NewAssetController() *AssetController {
	return &AssetController{
		service: services.NewAssetService(),
	}
}

func (c *AssetController) CreateAsset(ctx *gin.Context) {
	var req dtos.CreateAssetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	asset, err := c.service.CreateAsset(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, asset)
}

func (c *AssetController) GetAssets(ctx *gin.Context) {
	assets, total, err := c.service.GetAssets()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": assets, "total": total})
}

func (c *AssetController) GetAsset(ctx *gin.Context) {
	asset, err := c.service.GetAsset(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Asset not found"})
		return
	}
	ctx.JSON(http.StatusOK, asset)
}

func (c *AssetController) UpdateAsset(ctx *gin.Context) {
	var req dtos.UpdateAssetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateAsset(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Asset updated successfully"})
}

func (c *AssetController) DeleteAsset(ctx *gin.Context) {
	if err := c.service.DeleteAsset(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Asset deleted successfully"})
}