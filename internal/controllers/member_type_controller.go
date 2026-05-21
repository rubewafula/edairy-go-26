package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
	"gorm.io/gorm"
)

// MemberTypeController handles HTTP requests for member type management.
type MemberTypeController struct {
	service *services.MemberTypeService
}

// NewMemberTypeController creates a new instance of MemberTypeController.
func NewMemberTypeController() *MemberTypeController {
	return &MemberTypeController{
		service: services.NewMemberTypeService(),
	}
}

// CreateMemberType handles the creation of a new member type.
// @Summary Create a new member type
// @Accept json
// @Produce json
// @Param memberType body dtos.CreateMemberTypeRequest true "Member Type data"
// @Success 201 {object} dtos.MemberTypeResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 422 {object} map[string]string "Validation Error"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /member-types [post]
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

// GetMemberTypes retrieves a paginated list of member types.
// @Summary Get all member types
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} map[string]interface{} "data: []dtos.MemberTypeResponse, total: int64"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /member-types [get]
func (c *MemberTypeController) GetMemberTypes(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	results, total, err := c.service.GetMemberTypes(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

// GetMemberType retrieves a single member type by ID.
// @Summary Get a member type by ID
// @Produce json
// @Param id path string true "Member Type ID"
// @Success 200 {object} dtos.MemberTypeResponse
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /member-types/{id} [get]
func (c *MemberTypeController) GetMemberType(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.GetMemberType(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Member type not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// UpdateMemberType updates an existing member type.
// @Summary Update a member type
// @Accept json
// @Produce json
// @Param id path string true "Member Type ID"
// @Param memberType body dtos.UpdateMemberTypeRequest true "Updated Member Type data"
// @Success 200 {object} map[string]string "message: Member type updated successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 422 {object} map[string]string "Validation Error"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /member-types/{id} [put]
func (c *MemberTypeController) UpdateMemberType(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateMemberTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateMemberType(id, req); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Member type not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Member type updated successfully"})
}

// DeleteMemberType soft deletes a member type.
// @Summary Delete a member type
// @Produce json
// @Param id path string true "Member Type ID"
// @Success 200 {object} map[string]string "message: Member type deleted successfully"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /member-types/{id} [delete]
func (c *MemberTypeController) DeleteMemberType(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.DeleteMemberType(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Member type not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Member type deleted successfully"})
}
