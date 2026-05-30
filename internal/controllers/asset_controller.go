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
		log.Printf("[AssetController.CreateAsset] JSON binding error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[AssetController.CreateAsset] Validation error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	asset, err := c.service.CreateAsset(req)
	if err != nil {
		log.Printf("[AssetController.CreateAsset] Service error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, asset)
}

func (c *AssetController) GetAssets(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	assetCode := ctx.Query("asset_code")
	assetName := ctx.Query("asset_name")
	categoryID := ctx.Query("asset_category_id")
	fromDate := ctx.Query("asset_aquisition_date_from")
	toDate := ctx.Query("asset_aquisition_date_to")

	assets, total, err := c.service.GetAssets(page, limit, assetCode, assetName, categoryID, fromDate, toDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": assets, "total": total})
}

func (c *AssetController) GetAsset(ctx *gin.Context) {
	asset, err := c.service.GetAsset(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}
	ctx.JSON(http.StatusOK, asset)
}

func (c *AssetController) UpdateAsset(ctx *gin.Context) {
	var req dtos.UpdateAssetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateAsset(ctx.Param("id"), req); err != nil {
		log.Printf("[AssetController.UpdateAsset] Error updating asset %s: %v", ctx.Param("id"), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Asset updated successfully"})
}

func (c *AssetController) DeleteAsset(ctx *gin.Context) {
	if err := c.service.DeleteAsset(ctx.Param("id")); err != nil {
		log.Printf("[AssetController.DeleteAsset] Error deleting asset %s: %v", ctx.Param("id"), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Asset deleted successfully"})
}

func (c *AssetController) ImportAssets(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.ImportAssets(file, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": "Asset import started in the background. Check logs for status."})
}
