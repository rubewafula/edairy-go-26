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

func (c *MemberPayrollDeductionController) parseFilters(ctx *gin.Context) services.MemberPayrollDeductionFilters {
	return services.MemberPayrollDeductionFilters{
		MemberID:        firstNonEmptyQuery(ctx, "member_id", "memberID"),
		MemberNo:        firstNonEmptyQuery(ctx, "member_no", "memberNo"),
		PayrollID:       firstNonEmptyQuery(ctx, "member_payroll_id", "payroll_id", "payrollID"),
		DeductionTypeID: firstNonEmptyQuery(ctx, "deduction_type_id", "deductionTypeID"),
		DeductionMonth:  firstNonEmptyQuery(ctx, "deduction_month", "deductionMonth"),
		FiscalYear:      firstNonEmptyQuery(ctx, "fiscal_year", "fiscalYear"),
		Settled:         firstNonEmptyQuery(ctx, "settled"),
		Confirmed:       firstNonEmptyQuery(ctx, "confirmed"),
		DateFrom:        firstNonEmptyQuery(ctx, "date_from", "dateFrom"),
		DateTo:          firstNonEmptyQuery(ctx, "date_to", "dateTo"),
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
// @Param member_id query string false "Filter by member ID"
// @Param member_no query string false "Filter by member number"
// @Param member_payroll_id query string false "Filter by payroll ID"
// @Param deduction_type_id query string false "Filter by deduction type ID"
// @Param date_from query string false "Filter by transaction date from (YYYY-MM-DD)"
// @Param date_to query string false "Filter by transaction date to (YYYY-MM-DD)"
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
	filters := c.parseFilters(ctx)

	deductions, total, err := c.service.GetMemberPayrollDeductions(page, limit, filters)
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
	filters := c.parseFilters(ctx)
	format := ctx.DefaultQuery("format", "csv")

	userID := ctx.GetUint64("user_id")
	if err := c.service.ExportDeductions(userID, filters, format); err != nil {
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
