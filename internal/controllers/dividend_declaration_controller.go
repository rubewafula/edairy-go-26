package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type DividendDeclarationController struct {
	service *services.DividendDeclarationService
}

func NewDividendDeclarationController() *DividendDeclarationController {
	return &DividendDeclarationController{
		service: services.NewDividendDeclarationService(),
	}
}

func (c *DividendDeclarationController) CreateDeclaration(ctx *gin.Context) {
	var req dtos.CreateDividendDeclarationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	declaration, err := c.service.CreateDeclaration(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, declaration)
}

func (c *DividendDeclarationController) GetDeclarations(ctx *gin.Context) {
	declarations, total, err := c.service.GetDeclarations()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Data": declarations, "Total": total})
}

func (c *DividendDeclarationController) GetDeclaration(ctx *gin.Context) {
	declaration, err := c.service.GetDeclaration(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Dividend declaration not found"})
		return
	}
	ctx.JSON(http.StatusOK, declaration)
}

func (c *DividendDeclarationController) UpdateDeclaration(ctx *gin.Context) {
	var req dtos.UpdateDividendDeclarationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateDeclaration(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Dividend declaration updated successfully"})
}

func (c *DividendDeclarationController) DeleteDeclaration(ctx *gin.Context) {
	if err := c.service.DeleteDeclaration(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Dividend declaration deleted successfully"})
}
