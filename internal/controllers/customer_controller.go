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
)

type CustomerController struct {
	service *services.CustomerService
}

func NewCustomerController() *CustomerController {
	return &CustomerController{
		service: services.NewCustomerService(),
	}
}

func (c *CustomerController) CreateCustomer(ctx *gin.Context) {
	var req dtos.CreateCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	customer, err := c.service.CreateCustomer(req, userID)
	if err != nil {
		log.Printf("[CustomerController.CreateCustomer] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, customer)
}

func (c *CustomerController) GetCustomers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	customers, total, err := c.service.GetCustomers(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": customers, "total": total})
}

func (c *CustomerController) GetCustomer(ctx *gin.Context) {
	customer, err := c.service.GetCustomer(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

func (c *CustomerController) UpdateCustomer(ctx *gin.Context) {
	var req dtos.UpdateCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateCustomer(ctx.Param("id"), req, userID); err != nil {
		log.Printf("[CustomerController.UpdateCustomer] Error updating customer %s: %v", ctx.Param("id"), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully"})
}

func (c *CustomerController) DeleteCustomer(ctx *gin.Context) {
	if err := c.service.DeleteCustomer(ctx.Param("id")); err != nil {
		log.Printf("[CustomerController.DeleteCustomer] Error deleting customer %s: %v", ctx.Param("id"), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

func (c *CustomerController) ImportCustomers(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.ImportCustomers(file, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": "Import started in background"})
}

// ExportCustomers triggers the background generation of a customer export.
func (c *CustomerController) ExportCustomers(ctx *gin.Context) {
	status := ctx.Query("status")
	format := ctx.Query("format")
	userID := ctx.GetUint64("user_id")

	if err := c.service.ExportCustomers(userID, status, format); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "Export initiated. Check your notifications for the download link shortly.",
	})
}

// DownloadExportFile serves the generated CSV or PDF file for download.
func (c *CustomerController) DownloadExportFile(ctx *gin.Context) {
	filename := ctx.Param("filename")
	safeFilename := filepath.Base(filename)
	filePath := filepath.Join("./storage/exports", safeFilename)

	ctx.File(filePath)
}

// GetImportErrors returns the validation/processing errors for a specific import ID.
func (c *CustomerController) GetImportErrors(ctx *gin.Context) {
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
