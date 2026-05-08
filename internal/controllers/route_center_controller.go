package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type RouteCenterController struct {
	service *services.RouteCenterService
}

func NewRouteCenterController() *RouteCenterController {
	return &RouteCenterController{
		service: services.NewRouteCenterService(),
	}
}

func (c *RouteCenterController) CreateCenter(ctx *gin.Context) {
	var req dtos.CreateRouteCenterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	center, err := c.service.CreateCenter(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, center)
}

func (c *RouteCenterController) GetCenters(ctx *gin.Context) {
	centers, total, err := c.service.GetCenters()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": centers, "total": total})
}

func (c *RouteCenterController) GetCenter(ctx *gin.Context) {
	center, err := c.service.GetCenter(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Route center not found"})
		return
	}
	ctx.JSON(http.StatusOK, center)
}

func (c *RouteCenterController) UpdateCenter(ctx *gin.Context) {
	var req dtos.UpdateRouteCenterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateCenter(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Route center updated successfully"})
}

func (c *RouteCenterController) DeleteCenter(ctx *gin.Context) {
	if err := c.service.DeleteCenter(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Route center deleted successfully"})
}
