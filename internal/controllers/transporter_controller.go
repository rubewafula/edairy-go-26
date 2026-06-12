package controllers

import (
	"bytes"
	"errors"
	"io"
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

type TransporterController struct {
	service *services.TransporterService
}

func NewTransporterController() *TransporterController {
	return &TransporterController{
		service: services.NewTransporterService(),
	}
}

func (c *TransporterController) logRawRequest(ctx *gin.Context) {

	body, _ := io.ReadAll(ctx.Request.Body)

	log.Printf(`
	Method: %s
	URL: %s
	Headers: %+v
	Body: %s
	`,
		ctx.Request.Method,
		ctx.Request.URL.String(),
		ctx.Request.Header,
		string(body),
	)

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
}

func (c *TransporterController) CreateTransporter(ctx *gin.Context) {
	var req dtos.CreateTransporterRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Println("Error binding JSON: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Println("Error validating binding JSON: %s", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	transporter, err := c.service.CreateTransporter(req)
	if err != nil {
		log.Println("Error Creating transporter: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, transporter)
}

func (c *TransporterController) GetTransporters(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	transporters, total, err := c.service.GetTransporters(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": transporters, "total": total})
}

func (c *TransporterController) GetTransporter(ctx *gin.Context) {
	transporter, err := c.service.GetTransporter(ctx.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Transporter not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, transporter)
}

func (c *TransporterController) UpdateTransporter(ctx *gin.Context) {
	var req dtos.UpdateTransporterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateTransporter(ctx.Param("id"), req); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Transporter not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Transporter updated successfully"})
}

func (c *TransporterController) DeleteTransporter(ctx *gin.Context) {
	if err := c.service.DeleteTransporter(ctx.Param("id")); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Transporter not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// ImportTransporters handles the bulk upload of transporters.
func (c *TransporterController) ImportTransporters(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Excel or CSV file is required"})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.ImportTransporters(file, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "Transporter import started in the background. You will be notified upon completion.",
	})
}

// ExportTransporters triggers the background generation of a transporter CSV export.
func (c *TransporterController) ExportTransporters(ctx *gin.Context) {
	status := ctx.Query("status")
	format := ctx.Query("format")
	userID := ctx.GetUint64("user_id")

	if err := c.service.ExportTransporters(userID, status, format); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "Export initiated. Check your notifications for the download link shortly.",
	})
}

// DownloadExportFile serves the generated CSV file for download.
func (c *TransporterController) DownloadExportFile(ctx *gin.Context) {
	filename := ctx.Param("filename")
	safeFilename := filepath.Base(filename)
	filePath := filepath.Join("./storage/exports", safeFilename)

	ctx.File(filePath)
}

// GetImportErrors returns the validation/processing errors for a specific import ID.
func (c *TransporterController) GetImportErrors(ctx *gin.Context) {
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
