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

type CoolerMilkCollectionController struct {
	service *services.CoolerMilkCollectionService
}

func NewCoolerMilkCollectionController() *CoolerMilkCollectionController {
	return &CoolerMilkCollectionController{
		service: services.NewCoolerMilkCollectionService(),
	}
}

func (c *CoolerMilkCollectionController) CreateCollection(ctx *gin.Context) {
	var req dtos.CreateCoolerMilkCollectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	collection, err := c.service.CreateCollection(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, collection)
}

func (c *CoolerMilkCollectionController) GetCollections(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	collections, total, err := c.service.GetCollections(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": collections, "total": total})
}

func (c *CoolerMilkCollectionController) GetCollection(ctx *gin.Context) {
	collection, err := c.service.GetCollection(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Collection entry not found"})
		return
	}
	ctx.JSON(http.StatusOK, collection)
}

func (c *CoolerMilkCollectionController) UpdateCollection(ctx *gin.Context) {
	var req dtos.UpdateCoolerMilkCollectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateCollection(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Collection entry updated successfully"})
}

func (c *CoolerMilkCollectionController) DeleteCollection(ctx *gin.Context) {
	if err := c.service.DeleteCollection(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Collection entry deleted successfully"})
}
