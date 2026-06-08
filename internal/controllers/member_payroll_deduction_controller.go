package controllers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"gorm.io/gorm"
)

// MemberPayrollDeductionController handles requests related to member payroll deductions.
type MemberPayrollDeductionController struct {
	service *services.MemberPayrollDeductionService
}

// NewMemberPayrollDeductionController creates a new MemberPayrollDeductionController.
func NewMemberPayrollDeductionController() *MemberPayrollDeductionController {
	return &MemberPayrollDeductionController{
		service: services.NewMemberPayrollDeductionService(),
	}
}

// GetMemberPayrollDeductions godoc
// @Summary Get all member payroll deductions
// @Description Get a list of all member payroll deductions with optional filtering and pagination
// @Tags Member Payroll Deductions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param memberID query string false "Filter by member ID"
// @Param payrollID query string false "Filter by payroll ID"
// @Param deductionMonth query string false "Filter by deduction month"
// @Param fiscalYear query string false "Filter by fiscal year"
// @Param settled query string false "Filter by settled status (e.g., '0' or '1')"
// @Param confirmed query string false "Filter by confirmed status (e.g., '0' or '1')"
// @Success 200 {object} map[string]interface{} "List of member payroll deductions"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /member-payroll-deductions [get]
func (c *MemberPayrollDeductionController) GetMemberPayrollDeductions(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	memberID := ctx.Query("memberID")
	payrollID := ctx.Query("payrollID")
	deductionMonth := ctx.Query("deductionMonth")
	fiscalYear := ctx.Query("fiscalYear")
	settled := ctx.Query("settled")
	confirmed := ctx.Query("confirmed")

	deductions, total, err := c.service.GetMemberPayrollDeductions(page, limit, memberID, payrollID, deductionMonth, fiscalYear, settled, confirmed)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  deductions,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetMemberPayrollDeduction godoc
// @Summary Get a single member payroll deduction by ID
// @Description Get details of a specific member payroll deduction by its ID
// @Tags Member Payroll Deductions
// @Accept json
// @Produce json
// @Param id path string true "Member Payroll Deduction ID"
// @Success 200 {object} map[string]interface{} "Member payroll deduction details"
// @Failure 404 {object} map[string]interface{} "Member payroll deduction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /member-payroll-deductions/{id} [get]
func (c *MemberPayrollDeductionController) GetMemberPayrollDeduction(ctx *gin.Context) {
	id := ctx.Param("id")

	deduction, err := c.service.GetMemberPayrollDeduction(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Member payroll deduction not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": deduction})
}

func (c *MemberPayrollDeductionController) ExportDeductions(ctx *gin.Context) {
	memberID := ctx.Query("memberID")
	payrollID := ctx.Query("payrollID")
	deductionMonth := ctx.Query("deductionMonth")
	fiscalYear := ctx.Query("fiscalYear")
	settled := ctx.Query("settled")
	confirmed := ctx.Query("confirmed")
	format := ctx.DefaultQuery("format", "csv")

	userID := ctx.GetUint64("user_id")
	if err := c.service.ExportDeductions(userID, memberID, payrollID, deductionMonth, fiscalYear, settled, confirmed, format); err != nil {
		log.Printf("[MemberPayrollDeductionController.ExportDeductions] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"message": "Export started. You will receive a notification when it is ready."})
}

func (c *MemberPayrollDeductionController) DownloadExportFile(ctx *gin.Context) {
	filename := filepath.Base(ctx.Param("filename"))
	filePath := filepath.Join("./storage/exports", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	ctx.File(filePath)
}
