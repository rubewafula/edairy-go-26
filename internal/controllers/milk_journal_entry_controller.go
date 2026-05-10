package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type MilkJournalEntryController struct {
	service *services.MilkJournalEntryService
}

func NewMilkJournalEntryController() *MilkJournalEntryController {
	return &MilkJournalEntryController{
		service: services.NewMilkJournalEntryService(),
	}
}

func (c *MilkJournalEntryController) CreateEntry(ctx *gin.Context) {
	var req dtos.CreateMilkJournalEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	entry, err := c.service.CreateEntry(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, entry)
}

func (c *MilkJournalEntryController) GetEntries(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	entries, total, err := c.service.GetEntries(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Data": entries, "Total": total})
}

func (c *MilkJournalEntryController) GetEntry(ctx *gin.Context) {
	entry, err := c.service.GetEntry(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Entry not found"})
		return
	}
	ctx.JSON(http.StatusOK, entry)
}

func (c *MilkJournalEntryController) UpdateEntry(ctx *gin.Context) {
	var req dtos.UpdateMilkJournalEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateEntry(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Entry updated successfully"})
}

func (c *MilkJournalEntryController) DeleteEntry(ctx *gin.Context) {
	if err := c.service.DeleteEntry(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Entry deleted successfully"})
}

func (c *MilkJournalEntryController) GetStrayEntries(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	entries, total, err := c.service.GetStrayEntries(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Data": entries, "Total": total})
}

func (c *MilkJournalEntryController) UploadEntries(ctx *gin.Context) {
	// Implementation for Excel upload would go here
	ctx.JSON(http.StatusNotImplemented, gin.H{"Message": "Bulk upload via XLS not yet implemented"})
}
