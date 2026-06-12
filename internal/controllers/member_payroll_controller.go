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

type MemberPayrollController struct {
	service *services.MemberPayrollService
}

func NewMemberPayrollController() *MemberPayrollController {
	return &MemberPayrollController{
		service: services.NewMemberPayrollService(),
	}
}

func (c *MemberPayrollController) Create(ctx *gin.Context) {
	var req dtos.CreateMemberPayrollRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	res, err := c.service.Create(req, userID)
	if err != nil {
		log.Printf("[MemberPayrollController.Create] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *MemberPayrollController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	res, total, err := c.service.List(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *MemberPayrollController) Approve(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid payroll ID",
		})
		return
	}

	userID := ctx.GetUint64("user_id")

	// Define a local struct to capture the approval flag
	var req struct {
		IsApproved bool `json:"is_approved"`
	}

	// Bind the JSON body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "is_approved field is required"})
		return
	}
	log.Printf("Received payroll approve Request with status: %v", req.IsApproved)

	payroll, err := c.service.Approve(id, userID, req.IsApproved)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Member payroll not found"})
			return
		}
		log.Printf("[MemberPayrollController.Approve] Error processing payroll approval/rejection %d: %v", id, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message := "Payroll approval/rejection initiated successfully"
	ctx.JSON(http.StatusOK, gin.H{"message": message, "payroll": payroll})
}

func (c *MemberPayrollController) Confirm(ctx *gin.Context) {

	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Invalid payroll",
		})
		return
	}

	userID := ctx.GetUint64("user_id")

	payroll, err := c.service.Confirm(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Member payroll not found"})
			return
		}
		log.Printf("[MemberPayrollController.Confirm] Error confirming payroll %d: %v", id, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Member payroll confirmed successfully", "payroll": payroll})
}

func (c *MemberPayrollController) Get(ctx *gin.Context) {
	res, err := c.service.Get(ctx.Param("id"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Payroll record not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *MemberPayrollController) GetGenerationErrors(ctx *gin.Context) {
	payrollIDStr := ctx.Param("payrollID")
	payrollID, err := strconv.ParseUint(payrollIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid payroll ID",
		})
		return
	}

	errors, err := c.service.GetGenerationErrors(payrollID)
	if err != nil {
		log.Printf("[MemberPayrollController.GetGenerationErrors] Error fetching generation errors for payroll %d: %v", payrollID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payroll generation errors"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": errors})
}

func (c *MemberPayrollController) GetApprovalErrors(ctx *gin.Context) {
	payrollIDStr := ctx.Param("payrollID")
	payrollID, err := strconv.ParseUint(payrollIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payroll ID"})
		return
	}
	errors, err := c.service.GetApprovalErrors(payrollID)
	if err != nil {
		log.Printf("[MemberPayrollController.GetApprovalErrors] Error fetching approval errors for payroll %d: %v", payrollID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payroll approval errors"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": errors})
}
