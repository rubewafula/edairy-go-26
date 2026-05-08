package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type MemberTypeController struct {
	service *services.MemberTypeService
}

func NewMemberTypeController() *MemberTypeController {
	return &MemberTypeController{
		service: services.NewMemberTypeService(),
	}
}

func (c *MemberTypeController) CreateMemberType(ctx *gin.Context) {
	var req dtos.CreateMemberTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	memberType, err := c.service.CreateMemberType(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, memberType)
}

func (c *MemberTypeController) GetMemberTypes(ctx *gin.Context) {
	memberTypes, total, err := c.service.GetMemberTypes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": memberTypes, "total": total})
}

func (c *MemberTypeController) GetMemberType(ctx *gin.Context) {
	memberType, err := c.service.GetMemberType(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Member Type not found"})
		return
	}
	ctx.JSON(http.StatusOK, memberType)
}

func (c *MemberTypeController) UpdateMemberType(ctx *gin.Context) {
	var req dtos.UpdateMemberTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateMemberType(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Member Type updated successfully"})
}

func (c *MemberTypeController) DeleteMemberType(ctx *gin.Context) {
	if err := c.service.DeleteMemberType(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Member Type deleted successfully"})
}
