package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type SubRouteController struct {
	service *services.SubRouteService
}

func NewSubRouteController() *SubRouteController {
	return &SubRouteController{
		service: services.NewSubRouteService(),
	}
}

func (c *SubRouteController) CreateSubRoute(ctx *gin.Context) {
	var req dtos.CreateSubRouteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	subRoute, err := c.service.CreateSubRoute(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, subRoute)
}

func (c *SubRouteController) GetSubRoutes(ctx *gin.Context) {
	subRoutes, total, err := c.service.GetSubRoutes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": subRoutes, "total": total})
}

func (c *SubRouteController) GetSubRoute(ctx *gin.Context) {
	subRoute, err := c.service.GetSubRoute(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "SubRoute not found"})
		return
	}
	ctx.JSON(http.StatusOK, subRoute)
}

func (c *SubRouteController) UpdateSubRoute(ctx *gin.Context) {
	var req dtos.UpdateSubRouteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateSubRoute(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "SubRoute updated successfully"})
}

func (c *SubRouteController) DeleteSubRoute(ctx *gin.Context) {
	if err := c.service.DeleteSubRoute(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "SubRoute deleted successfully"})
}
