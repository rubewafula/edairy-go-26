package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type RouteController struct {
	service *services.RouteService
}

func NewRouteController() *RouteController {
	return &RouteController{
		service: services.NewRouteService(),
	}
}

func (c *RouteController) CreateRoute(ctx *gin.Context) {
	var req dtos.CreateRouteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	route, err := c.service.CreateRoute(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, route)
}

func (c *RouteController) GetRoutes(ctx *gin.Context) {
	routes, total, err := c.service.GetRoutes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": routes, "total": total})
}

func (c *RouteController) GetRoute(ctx *gin.Context) {
	route, err := c.service.GetRoute(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		return
	}
	ctx.JSON(http.StatusOK, route)
}

func (c *RouteController) UpdateRoute(ctx *gin.Context) {
	var req dtos.UpdateRouteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateRoute(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Route updated successfully"})
}

func (c *RouteController) DeleteRoute(ctx *gin.Context) {
	if err := c.service.DeleteRoute(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Route deleted successfully"})
}
