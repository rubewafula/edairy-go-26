package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type ExchangeVisitController struct {
	service *services.ExchangeVisitService
}

func NewExchangeVisitController() *ExchangeVisitController {
	return &ExchangeVisitController{
		service: services.NewExchangeVisitService(),
	}
}

func (c *ExchangeVisitController) CreateVisit(ctx *gin.Context) {
	var req dtos.CreateExchangeVisitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	visit, err := c.service.CreateVisit(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, visit)
}

func (c *ExchangeVisitController) GetVisits(ctx *gin.Context) {
	visits, total, err := c.service.GetVisits()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": visits, "total": total})
}

func (c *ExchangeVisitController) GetVisit(ctx *gin.Context) {
	visit, err := c.service.GetVisit(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Visit not found"})
		return
	}
	ctx.JSON(http.StatusOK, visit)
}

func (c *ExchangeVisitController) UpdateVisit(ctx *gin.Context) {
	var req dtos.UpdateExchangeVisitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateVisit(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Visit updated successfully"})
}

func (c *ExchangeVisitController) DeleteVisit(ctx *gin.Context) {
	if err := c.service.DeleteVisit(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Visit deleted successfully"})
}
