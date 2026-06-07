package controllers

import (
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

type EmployeeLeaveApplicationController struct {
	service *services.EmployeeLeaveApplicationService
}

func NewEmployeeLeaveApplicationController() *EmployeeLeaveApplicationController {
	return &EmployeeLeaveApplicationController{
		service: services.NewEmployeeLeaveApplicationService(),
	}
}

func (c *EmployeeLeaveApplicationController) Create(ctx *gin.Context) {
	var req dtos.CreateEmployeeLeaveApplicationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	res, err := c.service.CreateApplication(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (c *EmployeeLeaveApplicationController) List(ctx *gin.Context) {
	employeeID := ctx.Query("employee_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	results, total, err := c.service.GetApplications(employeeID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  results,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (c *EmployeeLeaveApplicationController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.service.GetApplication(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *EmployeeLeaveApplicationController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateEmployeeLeaveApplicationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateApplication(id, req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Application updated successfully"})
}

func (c *EmployeeLeaveApplicationController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetUint64("user_id")

	if err := c.service.DeleteApplication(id, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Application deleted successfully"})
}

func (c *EmployeeLeaveApplicationController) Export(ctx *gin.Context) {
	status := ctx.Query("status")
	format := ctx.DefaultQuery("format", "csv")
	userID := ctx.GetUint64("user_id")

	if err := c.service.ExportApplications(userID, status, format); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "Export initiated. Check your notifications for the download link shortly.",
	})
}

func (c *EmployeeLeaveApplicationController) DownloadExportFile(ctx *gin.Context) {
	filename := ctx.Param("filename")
	safeFilename := filepath.Base(filename)
	filePath := filepath.Join("./storage/exports", safeFilename)

	ctx.File(filePath)
}
