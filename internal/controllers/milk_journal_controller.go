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
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	journal, err := c.service.CreateMilkJournal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, journal)
}

func (c *MilkJournalController) GetMilkJournals(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	journals, total, err := c.service.GetMilkJournals(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Data": journals, "Total": total})
}

func (c *MilkJournalController) GetMilkJournal(ctx *gin.Context) {
	journal, err := c.service.GetMilkJournal(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Milk Journal not found"})
		return
	}
	ctx.JSON(http.StatusOK, journal)
}

func (c *MilkJournalController) UpdateMilkJournal(ctx *gin.Context) {
	var req dtos.UpdateMilkJournalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateMilkJournal(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Milk Journal updated successfully"})
}

func (c *MilkJournalController) DeleteMilkJournal(ctx *gin.Context) {
	if err := c.service.DeleteMilkJournal(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Milk Journal deleted successfully"})
}
