package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type AssetCategoryController struct {
	service *services.AssetCategoryService
}

func NewAssetCategoryController() *AssetCategoryController {
	return &AssetCategoryController{
		service: services.NewAssetCategoryService(),
	}
}

func (c *AssetCategoryController) CreateCategory(ctx *gin.Context) {
	var req dtos.CreateAssetCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	category, err := c.service.CreateCategory(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, category)
}

func (c *AssetCategoryController) GetCategories(ctx *gin.Context) {
	categories, total, err := c.service.GetCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": categories, "total": total})
}

func (c *AssetCategoryController) GetCategory(ctx *gin.Context) {
	category, err := c.service.GetCategory(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Asset category not found"})
		return
	}
	ctx.JSON(http.StatusOK, category)
}

func (c *AssetCategoryController) UpdateCategory(ctx *gin.Context) {
	var req dtos.UpdateAssetCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateCategory(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Asset category updated successfully"})
}

func (c *AssetCategoryController) DeleteCategory(ctx *gin.Context) {
	if err := c.service.DeleteCategory(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Asset category deleted successfully"})
}
