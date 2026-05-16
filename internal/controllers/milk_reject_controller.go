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

type MilkRejectController struct {
	service *services.MilkRejectService
}

func NewMilkRejectController() *MilkRejectController {
	return &MilkRejectController{
		service: services.NewMilkRejectService(),
	}
}

func (c *MilkRejectController) CreateReject(ctx *gin.Context) {
	var req dtos.CreateMilkRejectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	reject, err := c.service.CreateReject(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, reject)
}

func (c *MilkRejectController) GetRejects(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	rejects, total, err := c.service.GetRejects(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": rejects, "total": total})
}

func (c *MilkRejectController) GetReject(ctx *gin.Context) {
	reject, err := c.service.GetReject(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Reject entry not found"})
		return
	}
	ctx.JSON(http.StatusOK, reject)
}

func (c *MilkRejectController) DeleteReject(ctx *gin.Context) {
	if err := c.service.DeleteReject(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Reject entry deleted successfully"})
}
