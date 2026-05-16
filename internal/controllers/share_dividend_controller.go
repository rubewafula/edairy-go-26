package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type ShareDividendController struct {
	service *services.ShareDividendService
}

func NewShareDividendController() *ShareDividendController {
	return &ShareDividendController{
		service: services.NewShareDividendService(),
	}
}

func (c *ShareDividendController) CreateDividend(ctx *gin.Context) {
	var req dtos.CreateShareDividendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	dividend, err := c.service.CreateDividend(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, dividend)
}

func (c *ShareDividendController) GetDividends(ctx *gin.Context) {
	dividends, total, err := c.service.GetDividends()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": dividends, "total": total})
}

func (c *ShareDividendController) GetDividend(ctx *gin.Context) {
	dividend, err := c.service.GetDividend(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Share dividend not found"})
		return
	}
	ctx.JSON(http.StatusOK, dividend)
}

func (c *ShareDividendController) UpdateDividend(ctx *gin.Context) {
	var req dtos.UpdateShareDividendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateDividend(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share dividend updated successfully"})
}

func (c *ShareDividendController) DeleteDividend(ctx *gin.Context) {
	if err := c.service.DeleteDividend(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share dividend deleted successfully"})
}
