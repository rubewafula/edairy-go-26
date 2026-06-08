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

type MemberPayslipController struct {
	service *services.MemberPayrollService
}

func NewMemberPayslipController() *MemberPayslipController {
	return &MemberPayslipController{
		service: services.NewMemberPayrollService(),
	}
}

func (c *MemberPayslipController) GetPayslips(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	payrollID := ctx.Query("payroll_id")
	memberID := ctx.Query("member_id")
	routeID := ctx.Query("route_id")

	res, total, err := c.service.GetPayslips(payrollID, memberID, routeID, page, limit)
	if err != nil {
		log.Printf("[MemberPayslipController.GetPayslips] Error fetching payslips: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *MemberPayslipController) ExportPayslips(ctx *gin.Context) {
	payrollID := ctx.Query("payroll_id")
	memberID := ctx.Query("member_id")
	routeID := ctx.Query("route_id")
	format := ctx.DefaultQuery("format", "csv")

	userID := ctx.GetUint64("user_id")
	if err := c.service.ExportPayslips(userID, payrollID, memberID, routeID, format); err != nil {
		log.Printf("[MemberPayslipController.ExportPayslips] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "Payslip export started in the background. You will receive a notification when it's ready for download.",
	})
}

func (c *MemberPayslipController) ExportStatements(ctx *gin.Context) {
	payslipID := ctx.Param("payslip_id")
	memberID := ctx.Param("member_id")

	if payslipID == "" || memberID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "payslip_id and member_id are mandatory"})
		return
	}
	format := ctx.DefaultQuery("format", "pdf") // PDF is usually preferred for statements

	userID := ctx.GetUint64("user_id")
	if err := c.service.ExportPayslipStatements(userID, payslipID, memberID, format); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "Payslip statement export started. You will be notified when the report is ready.",
	})
}

func (c *MemberPayslipController) DownloadExportFile(ctx *gin.Context) {
	filename := filepath.Base(ctx.Param("filename"))
	filePath := filepath.Join("./storage/exports", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Export file not found"})
		return
	}

	ctx.File(filePath)
}

func (c *MemberPayslipController) GetPayslip(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := c.service.GetPayslip(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Member payslip not found"})
			return
		}
		log.Printf("[MemberPayslipController.GetPayslip] Error fetching payslip %s: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}
