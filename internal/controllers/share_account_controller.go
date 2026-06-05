package controllers

import (
	"net/http"
	"os"
	"path/filepath"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type ShareAccountController struct {
	service *services.ShareAccountService
}

func NewShareAccountController() *ShareAccountController {
	return &ShareAccountController{
		service: services.NewShareAccountService(),
	}
}

func (c *ShareAccountController) CreateAccount(ctx *gin.Context) {
	var req dtos.CreateShareAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	account, err := c.service.CreateAccount(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, account)
}

func (c *ShareAccountController) GetAccounts(ctx *gin.Context) {
	memberID := ctx.Query("member_id")
	shareTypeID := ctx.Query("share_type_id")
	accounts, total, err := c.service.GetShareAccounts(memberID, shareTypeID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": accounts, "total": total})
}

func (c *ShareAccountController) GetAccount(ctx *gin.Context) {
	account, err := c.service.GetShareAccount(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Share account not found"})
		return
	}
	ctx.JSON(http.StatusOK, account)
}

func (c *ShareAccountController) UpdateAccount(ctx *gin.Context) {
	var req dtos.UpdateShareAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateAccount(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share account updated successfully"})
}

func (c *ShareAccountController) DeleteAccount(ctx *gin.Context) {
	if err := c.service.DeleteAccount(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Share account deleted successfully"})
}

func (c *ShareAccountController) ImportAccounts(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.ImportAccounts(file, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": "Share account import started in the background. You will receive a notification when it's ready."})
}

func (c *ShareAccountController) GetImportErrors(ctx *gin.Context) {
	importIDStr := ctx.Param("importid")
	importID, err := strconv.ParseUint(importIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid import ID"})
		return
	}
	errors, err := c.service.GetImportErrors(importID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch import errors"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": errors})
}

func (c *ShareAccountController) ExportAccounts(ctx *gin.Context) {
	memberID := ctx.Query("member_id")
	shareTypeID := ctx.Query("share_type_id")
	status := ctx.Query("status")
	reportType := ctx.DefaultQuery("format", "csv")

	userID := ctx.GetUint64("user_id")
	if err := c.service.ExportAccounts(userID, memberID, shareTypeID, status, reportType); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": "Share account export started in the background."})
}

func (c *ShareAccountController) DownloadExportFile(ctx *gin.Context) {
	filename := filepath.Base(ctx.Param("filename"))
	filePath := filepath.Join("./storage/exports", filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	ctx.File(filePath)
}
