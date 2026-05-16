package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type AssetDepreciationController struct {
	service *services.AssetDepreciationService
}

func NewAssetDepreciationController() *AssetDepreciationController {
	return &AssetDepreciationController{
		service: services.NewAssetDepreciationService(),
	}
}

func (c *AssetDepreciationController) CreateEntry(ctx *gin.Context) {
	var req dtos.CreateAssetDepreciationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	entry, err := c.service.CreateEntry(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, entry)
}

func (c *AssetDepreciationController) GetEntries(ctx *gin.Context) {
	entries, total, err := c.service.GetEntries()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": entries, "total": total})
}

func (c *AssetDepreciationController) GetEntry(ctx *gin.Context) {
	entry, err := c.service.GetEntry(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Depreciation entry not found"})
		return
	}
	ctx.JSON(http.StatusOK, entry)
}

func (c *AssetDepreciationController) DeleteEntry(ctx *gin.Context) {
	if err := c.service.DeleteEntry(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Depreciation entry deleted successfully"})
}
