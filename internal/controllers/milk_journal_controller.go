package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type MilkJournalController struct {
	service *services.MilkJournalService
}

func NewMilkJournalController() *MilkJournalController {
	return &MilkJournalController{
		service: services.NewMilkJournalService(),
	}
}

func (c *MilkJournalController) CreateMilkJournal(ctx *gin.Context) {
	var req dtos.CreateMilkJournalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("MilkJournal controller failed to bind: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("MilkJournal controller failed to Validate: %s", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	journal, err := c.service.CreateMilkJournal(req)
	if err != nil {
		log.Printf("MilkJournal controller failed to create hournal : %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("MIlk journal create success:  %v", *journal)
	ctx.JSON(http.StatusCreated, gin.H{"data": *journal, "message": "Milk journal posted successfully"})
}

func (c *MilkJournalController) GetMilkJournals(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	journals, total, err := c.service.GetMilkJournals(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": journals, "total": total})
}

func (c *MilkJournalController) GetMilkJournal(ctx *gin.Context) {
	journal, err := c.service.GetMilkJournal(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Milk Journal not found"})
		return
	}
	ctx.JSON(http.StatusOK, journal)
}

func (c *MilkJournalController) UpdateMilkJournal(ctx *gin.Context) {
	var req dtos.UpdateMilkJournalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateMilkJournal(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Milk Journal updated successfully"})
}

func (c *MilkJournalController) DeleteMilkJournal(ctx *gin.Context) {
	if err := c.service.DeleteMilkJournal(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Milk Journal deleted successfully"})
}

func (c *MilkJournalController) GetDailyJournals(ctx *gin.Context) {
	journals, err := c.service.GetDailyJournals()
	if err != nil {
		log.Printf("MilkJournal controller failed to get daily journals: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": journals})
}

func (c *MilkJournalController) ImportJournals(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.ImportJournals(file, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": "Milk journal import started in the background. Check logs for status."})
}

func (c *MilkJournalController) GetMilkJournalImportErrors(ctx *gin.Context) {
	importIDStr := ctx.Param("importid")
	importID, _ := strconv.ParseUint(importIDStr, 10, 64)

	errors, err := c.service.GetImportErrors(importID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch import errors"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": errors})
}
