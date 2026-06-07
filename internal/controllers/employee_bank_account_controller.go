package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type EmployeeBankAccountController struct {
	service *services.EmployeeBankAccountService
}

func NewEmployeeBankAccountController() *EmployeeBankAccountController {
	return &EmployeeBankAccountController{
		service: services.NewEmployeeBankAccountService(),
	}
}

func (c *EmployeeBankAccountController) Create(ctx *gin.Context) {
	var req dtos.CreateEmployeeBankAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint64("user_id") // Assuming user_id is set by auth middleware
	res, err := c.service.CreateAccount(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (c *EmployeeBankAccountController) List(ctx *gin.Context) {
	employeeID := ctx.Query("employee_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	results, total, err := c.service.GetAccounts(employeeID, page, limit)
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

func (c *EmployeeBankAccountController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.service.GetAccount(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *EmployeeBankAccountController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateEmployeeBankAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateAccount(id, req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Account updated successfully"})
}

func (c *EmployeeBankAccountController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetUint64("user_id")

	if err := c.service.DeleteAccount(id, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

// ImportAccounts handles the bulk upload of employee bank accounts.
func (c *EmployeeBankAccountController) ImportAccounts(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Excel or CSV file is required"})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.ImportAccounts(file, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "Import started in the background. You will be notified upon completion.",
	})
}

// GetImportErrors returns the validation/processing errors for a specific import ID.
func (c *EmployeeBankAccountController) GetImportErrors(ctx *gin.Context) {
	idStr := ctx.Param("id")
	importID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid import session ID"})
		return
	}

	errors, err := c.service.GetImportErrors(importID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, errors)
}

// ExportAccounts triggers the background generation of an employee bank account export.
func (c *EmployeeBankAccountController) ExportAccounts(ctx *gin.Context) {
	format := ctx.DefaultQuery("format", "csv")
	userID := ctx.GetUint64("user_id")

	if err := c.service.ExportAccounts(userID, format); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "Export initiated. Check your notifications for the download link shortly.",
	})
}

// DownloadExportFile serves the generated CSV or PDF file for download.
func (c *EmployeeBankAccountController) DownloadExportFile(ctx *gin.Context) {
	filename := ctx.Param("filename")
	safeFilename := filepath.Base(filename)
	filePath := filepath.Join("./storage/exports", safeFilename)

	ctx.File(filePath)
}
