package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type ShareTypeController struct {
	service *services.ShareTypeService
}

func NewShareTypeController() *ShareTypeController {
	return &ShareTypeController{
		service: services.NewShareTypeService(),
	}
}

func (c *ShareTypeController) CreateShareType(ctx *gin.Context) {
	var req dtos.CreateShareTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	shareType, err := c.service.CreateShareType(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, shareType)
}

func (c *ShareTypeController) GetShareTypes(ctx *gin.Context) {
	shareTypes, total, err := c.service.GetShareTypes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": shareTypes, "total": total})
}

func (c *ShareTypeController) GetShareType(ctx *gin.Context) {
	shareType, err := c.service.GetShareType(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Share type not found"})
		return
	}
	ctx.JSON(http.StatusOK, shareType)
}

func (c *ShareTypeController) UpdateShareType(ctx *gin.Context) {
	var req dtos.UpdateShareTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateShareType(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share type updated successfully"})
}

func (c *ShareTypeController) DeleteShareType(ctx *gin.Context) {
	if err := c.service.DeleteShareType(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share type deleted successfully"})
}
