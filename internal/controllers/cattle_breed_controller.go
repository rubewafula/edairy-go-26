package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type CattleBreedController struct {
	service *services.CattleBreedService
}

func NewCattleBreedController() *CattleBreedController {
	return &CattleBreedController{
		service: services.NewCattleBreedService(),
	}
}

func (c *CattleBreedController) CreateCattleBreed(ctx *gin.Context) {
	var req dtos.CreateCattleBreedRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	cattleBreed, err := c.service.CreateCattleBreed(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, cattleBreed)
}

func (c *CattleBreedController) GetCattleBreeds(ctx *gin.Context) {
	cattleBreeds, total, err := c.service.GetCattleBreeds()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": cattleBreeds, "total": total})
}

func (c *CattleBreedController) GetCattleBreed(ctx *gin.Context) {
	cattleBreed, err := c.service.GetCattleBreed(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cattle Breed not found"})
		return
	}
	ctx.JSON(http.StatusOK, cattleBreed)
}

func (c *CattleBreedController) UpdateCattleBreed(ctx *gin.Context) {
	var req dtos.UpdateCattleBreedRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateCattleBreed(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Cattle Breed updated successfully"})
}

func (c *CattleBreedController) DeleteCattleBreed(ctx *gin.Context) {
	if err := c.service.DeleteCattleBreed(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Cattle Breed deleted successfully"})
}
