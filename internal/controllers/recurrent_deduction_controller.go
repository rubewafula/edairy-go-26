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

// RecurrentDeductionController handles requests related to recurrent deductions.
type RecurrentDeductionController struct {
	recurrentDeductionService *services.RecurrentDeductionService
}

// NewRecurrentDeductionController creates a new RecurrentDeductionController.
func NewRecurrentDeductionController() *RecurrentDeductionController {
	return &RecurrentDeductionController{
		recurrentDeductionService: services.NewRecurrentDeductionService(),
	}
}

// GetRecurrentDeductions godoc
// @Summary Get all recurrent deductions
// @Description Get a list of all recurrent deductions with optional filtering and pagination
// @Tags Recurrent Deductions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param customerID query string false "Filter by customer ID"
// @Param customerType query string false "Filter by customer type (e.g., member, transporter)"
// @Param settled query string false "Filter by settled status (0 for unsettled, 1 for settled)"
// @Success 200 {object} map[string]interface{} "data: []dtos.RecurrentDeductionResponse"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /member-deductions [get]
func (c *RecurrentDeductionController) GetRecurrentDeductions(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	customerID := ctx.Query("customerID")
	customerType := ctx.Query("customerType")
	settled := ctx.Query("settled")

	deductions, total, err := c.recurrentDeductionService.GetRecurrentDeductions(page, limit, customerID, customerType, settled)
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

// GetRecurrentDeduction godoc
// @Summary Get a single recurrent deduction by ID
// @Description Get details of a specific recurrent deduction by its ID
// @Tags Recurrent Deductions
// @Accept json
// @Produce json
// @Param id path string true "Recurrent Deduction ID"
// @Success 200 {object} dtos.RecurrentDeductionResponse
// @Failure 404 {object} map[string]interface{} "Recurrent deduction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /member-deductions/{id} [get]
func (c *RecurrentDeductionController) GetRecurrentDeduction(ctx *gin.Context) {
	id := ctx.Param("id")

	deduction, err := c.recurrentDeductionService.GetRecurrentDeduction(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Recurrent deduction not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": deduction})
}

// ExportRecurrentDeductions godoc
// @Summary Export recurrent deductions
// @Description Export a list of recurrent deductions to CSV or PDF with optional filtering
// @Tags Recurrent Deductions
// @Accept json
// @Produce json
// @Param customerID query string false "Filter by customer ID"
// @Param customerType query string false "Filter by customer type (e.g., member, transporter)"
// @Param settled query string false "Filter by settled status (0 for unsettled, 1 for settled)"
// @Param format query string false "Export format (csv or pdf)" default(csv)
// @Success 202 {object} map[string]interface{} "message: Export started in the background"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /member-deductions/export [get]
func (c *RecurrentDeductionController) ExportRecurrentDeductions(ctx *gin.Context) {
	customerID := ctx.Query("customerID")
	customerType := ctx.Query("customerType")
	settled := ctx.Query("settled")
	format := ctx.DefaultQuery("format", "csv")

	userID := ctx.GetUint64("user_id")
	if err := c.recurrentDeductionService.ExportRecurrentDeductions(userID, customerID, customerType, settled, format); err != nil {
		log.Printf("[RecurrentDeductionController.ExportRecurrentDeductions] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "Recurrent deductions export started in the background. You will receive a notification when it's ready for download.",
	})
}

// DownloadExportFile godoc
// @Summary Download an exported recurrent deductions file
// @Description Download a previously generated recurrent deductions export file
// @Tags Recurrent Deductions
// @Produce application/octet-stream
// @Param filename path string true "Name of the file to download"
// @Success 200 {file} file "Exported file"
// @Failure 404 {object} map[string]interface{} "File not found"
// @Router /member-deductions/export/download/{filename} [get]
func (c *RecurrentDeductionController) DownloadExportFile(ctx *gin.Context) {
	filename := filepath.Base(ctx.Param("filename"))
	filePath := filepath.Join("./storage/exports", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Export file not found"})
		return
	}

	ctx.File(filePath)
}
