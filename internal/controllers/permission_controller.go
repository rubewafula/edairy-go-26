package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type PermissionController struct {
	service *services.PermissionService
}

func NewPermissionController() *PermissionController {
	return &PermissionController{
		service: services.NewPermissionService(),
	}
}

func (c *PermissionController) CreatePermission(ctx *gin.Context) {
	var req dtos.CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	permission, err := c.service.CreatePermission(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, permission)
}

func (c *PermissionController) GetPermissions(ctx *gin.Context) {
	permissions, total, err := c.service.GetPermissions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": permissions, "total": total})
}

func (c *PermissionController) GetPermission(ctx *gin.Context) {
	permission, err := c.service.GetPermission(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
		return
	}
	ctx.JSON(http.StatusOK, permission)
}

func (c *PermissionController) UpdatePermission(ctx *gin.Context) {
	var req dtos.UpdatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdatePermission(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Permission updated successfully"})
}

func (c *PermissionController) DeletePermission(ctx *gin.Context) {
	if err := c.service.DeletePermission(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Permission deleted successfully"})
}
