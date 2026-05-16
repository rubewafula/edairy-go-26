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

type EmployeeLeaveTypeController struct {
	service *services.EmployeeLeaveTypeService
}

func NewEmployeeLeaveTypeController() *EmployeeLeaveTypeController {
	return &EmployeeLeaveTypeController{
		service: services.NewEmployeeLeaveTypeService(),
	}
}

func (c *EmployeeLeaveTypeController) CreateEmployeeLeaveType(ctx *gin.Context) {
	var req dtos.CreateEmployeeLeaveTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}
	userID := ctx.GetUint64("user_id")
	leaveType, err := c.service.CreateEmployeeLeaveType(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, leaveType)
}

func (c *EmployeeLeaveTypeController) GetEmployeeLeaveTypes(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))
	results, total, err := c.service.GetEmployeeLeaveTypes(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *EmployeeLeaveTypeController) GetEmployeeLeaveType(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.GetEmployeeLeaveType(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound || result == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee leave type not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *EmployeeLeaveTypeController) UpdateEmployeeLeaveType(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateEmployeeLeaveTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateEmployeeLeaveType(id, req, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee leave type not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Employee leave type updated successfully"})
}

func (c *EmployeeLeaveTypeController) DeleteEmployeeLeaveType(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteEmployeeLeaveType(id, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee leave type not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Employee leave type deleted successfully"})
}
