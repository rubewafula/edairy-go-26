package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
	"gorm.io/gorm"
)

type TransporterPayrollController struct {
	service *services.TransporterPayrollService
}

func NewTransporterPayrollController() *TransporterPayrollController {
	return &TransporterPayrollController{
		service: services.NewTransporterPayrollService(),
	}
}

func (c *TransporterPayrollController) Create(ctx *gin.Context) {
	var req dtos.CreateTransporterPayrollRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[TransporterPayrollController.Create] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[TransporterPayrollController.Create] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	res, err := c.service.Create(req, userID)
	if err != nil {
		log.Printf("[TransporterPayrollController.Create] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transporter payroll"})
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *TransporterPayrollController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	res, total, err := c.service.List(page, limit)
	if err != nil {
		log.Printf("[TransporterPayrollController.List] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payroll list"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *TransporterPayrollController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.service.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Transporter payroll not found"})
			return
		}
		log.Printf("[TransporterPayrollController.Get] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payroll details"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *TransporterPayrollController) Confirm(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payroll ID"})
		return
	}

	userID := ctx.GetUint64("user_id")
	payroll, err := c.service.Confirm(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Transporter payroll not found"})
			return
		}
		log.Printf("[TransporterPayrollController.Confirm] Error confirming payroll %d: %v", id, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Confirmation failed. Ensure payroll is in draft status."})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Transporter payroll confirmed successfully", "payroll": payroll})
}

func (c *TransporterPayrollController) Approve(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payroll ID"})
		return
	}

	userID := ctx.GetUint64("user_id")
	payroll, err := c.service.Approve(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Transporter payroll not found"})
			return
		}
		log.Printf("[TransporterPayrollController.Approve] Error approving payroll %d: %v", id, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Approval failed. Ensure payroll is confirmed and has no errors."})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Transporter payroll approval initiated", "payroll": payroll})
}
