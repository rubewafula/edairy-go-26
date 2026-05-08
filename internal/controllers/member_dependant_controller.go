package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type MemberDependantController struct {
	service *services.MemberDependantService
}

func NewMemberDependantController() *MemberDependantController {
	return &MemberDependantController{
		service: services.NewMemberDependantService(),
	}
}

func (c *MemberDependantController) CreateDependant(ctx *gin.Context) {
	var req dtos.CreateMemberDependantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	dependant, err := c.service.CreateMemberDependant(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, dependant)
}

func (c *MemberDependantController) GetDependants(ctx *gin.Context) {
	dependants, total, err := c.service.GetMemberDependants()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": dependants, "total": total})
}

func (c *MemberDependantController) GetDependant(ctx *gin.Context) {
	dependant, err := c.service.GetMemberDependant(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Dependant not found"})
		return
	}
	ctx.JSON(http.StatusOK, dependant)
}

func (c *MemberDependantController) UpdateDependant(ctx *gin.Context) {
	var req dtos.UpdateMemberDependantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateMemberDependant(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Dependant updated successfully"})
}

func (c *MemberDependantController) DeleteDependant(ctx *gin.Context) {
	if err := c.service.DeleteMemberDependant(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Dependant deleted successfully"})
}
