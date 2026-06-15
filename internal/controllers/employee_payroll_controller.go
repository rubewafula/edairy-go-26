package controllers

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
	"gorm.io/gorm"
)

type EmployeePayrollController struct {
	service *services.EmployeePayrollService
}

func NewEmployeePayrollController() *EmployeePayrollController {
	return &EmployeePayrollController{
		service: services.NewEmployeePayrollService(),
	}
}

// CreatePayroll godoc
// @Summary Initiate employee payroll generation
// @Description Triggers the background process to generate employee payrolls for all pending pay date ranges.
// @Tags Employee Payrolls
// @Accept json
// @Produce json
// @Success 202 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employee-payrolls [post]
func (c *EmployeePayrollController) CreatePayroll(ctx *gin.Context) {
	var req dtos.CreateEmployeePayrollRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")

	if err := c.service.CreatePayroll(req, userID); err != nil {
		log.Printf("[EmployeePayrollController.CreatePayroll] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"message": "Employee payroll generation started in the background"})
}

// ConfirmPayroll godoc
// @Summary Confirm an employee payroll
// @Description Confirms a generated employee payroll, moving it from draft/incomplete to confirmed status.
// @Tags Employee Payrolls
// @Accept json
// @Produce json
// @Param id path string true "Payroll ID"
// @Success 200 {object} models.EmployeePayroll
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employee-payrolls/{id}/confirm [put]
func (c *EmployeePayrollController) ConfirmPayroll(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payroll ID"})
		return
	}

	userID := ctx.GetUint64("user_id")
	payroll, err := c.service.ConfirmPayroll(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee payroll not found"})
			return
		}
		log.Printf("[EmployeePayrollController.ConfirmPayroll] Error confirming payroll %d: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, payroll)
}

// ApprovePayroll godoc
// @Summary Approve an employee payroll
// @Description Approves a confirmed employee payroll, triggering GL postings in the background.
// @Tags Employee Payrolls
// @Accept json
// @Produce json
// @Param id path string true "Payroll ID"
// @Success 202 {object} models.EmployeePayroll
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employee-payrolls/{id}/approve [put]
func (c *EmployeePayrollController) ApprovePayroll(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payroll ID"})
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

	payroll, err := c.service.ApprovePayroll(id, userID, req.IsApproved)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee payroll not found"})
			return
		}
		log.Printf("[EmployeePayrollController.ApprovePayroll] Error processing payroll approval/rejection %d: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	message := "Payroll approval/rejection initiated successfully"
	ctx.JSON(http.StatusAccepted, gin.H{"message": message, "data": payroll})
}

// GetEmployeePayrolls godoc
// @Summary Get all employee payrolls
// @Description Retrieve a list of all employee payrolls
// @Tags Employee Payrolls
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {array} dtos.EmployeePayrollResponse
// @Failure 500 {object} map[string]string
// @Router /employee-payrolls [get]
func (c *EmployeePayrollController) GetEmployeePayrolls(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	results, total, err := c.service.GetEmployeePayrolls(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

// GetEmployeePayroll godoc
// @Summary Get a single employee payroll by ID
// @Description Retrieve an employee payroll by its ID
// @Tags Employee Payrolls
// @Produce json
// @Param id path string true "Payroll ID"
// @Success 200 {object} dtos.EmployeePayrollResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employee-payrolls/{id} [get]
func (c *EmployeePayrollController) GetEmployeePayroll(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.GetEmployeePayroll(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee payroll not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// DeleteEmployeePayroll godoc
// @Summary Delete an employee payroll
// @Description Delete an employee payroll by its ID. Only draft or incomplete payrolls can be deleted.
// @Tags Employee Payrolls
// @Produce json
// @Param id path string true "Payroll ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employee-payrolls/{id} [delete]
func (c *EmployeePayrollController) DeleteEmployeePayroll(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payroll ID"})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteEmployeePayroll(utils.Uint64ToString(id), userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee payroll not found"})
			return
		}
		log.Printf("[EmployeePayrollController.DeleteEmployeePayroll] Error deleting payroll %d: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Employee payroll deleted successfully"})
}

// ExportStatement godoc
// @Summary Export employee payslip statement
// @Description Initiates background generation of a detailed payslip statement (PDF/CSV)
// @Tags Employee Payrolls
// @Produce json
// @Param employee_id path string true "Employee ID"
// @Param payroll_id path string true "Payroll ID"
// @Param report_type query string false "Report type (pdf/csv)" default(pdf)
// @Success 202 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employee-payslips/statements/{employee_id}/{payroll_id} [get]
func (c *EmployeePayrollController) ExportStatement(ctx *gin.Context) {
	employeeID := ctx.Param("employee_id")
	payrollID := ctx.Param("payroll_id")
	reportType := ctx.DefaultQuery("report_type", "pdf")
	userID := ctx.GetUint64("user_id")

	if err := c.service.ExportPayslipStatement(userID, employeeID, payrollID, reportType); err != nil {
		log.Printf("[EmployeePayrollController.ExportStatement] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate payslip export"})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": "Payslip statement generation started. You will be notified when the download is ready."})
}

func (c *EmployeePayrollController) DownloadExportFile(ctx *gin.Context) {
	filename := ctx.Param("filename")
	filePath := filepath.Join("./storage/exports", filename)
	ctx.File(filePath)
}

// ExportPayslips godoc
// @Summary Export employee payslips list
// @Description Initiates background generation of a list of employee payslips (PDF/CSV) based on filters
// @Tags Employee Payrolls
// @Produce json
// @Param payroll_id query string false "Payroll ID filter"
// @Param payroll_month query string false "Month filter"
// @Param payroll_year query string false "Year filter"
// @Param format query string false "Report format (pdf/csv)" default(csv)
// @Success 202 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employee-payslips/export [get]
func (c *EmployeePayrollController) ExportPayslips(ctx *gin.Context) {
	format := ctx.DefaultQuery("format", "csv")
	userID := ctx.GetUint64("user_id")

	filters := make(map[string]string)
	if ctx.Query("payroll_id") != "" {
		filters["payroll_id"] = ctx.Query("payroll_id")
	}
	if ctx.Query("payroll_month") != "" {
		filters["payroll_month"] = ctx.Query("payroll_month")
	}
	if ctx.Query("payroll_year") != "" {
		filters["payroll_year"] = ctx.Query("payroll_year")
	}

	if err := c.service.ExportPayslips(userID, filters, format); err != nil {
		log.Printf("[EmployeePayrollController.ExportPayslips] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate payslips export"})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": "Employee payslips export started. You will be notified when the download is ready."})
}
