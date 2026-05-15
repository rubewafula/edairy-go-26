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

type SupplyRejectController struct {
	service *services.SupplyRejectService
}

func NewSupplyRejectController() *SupplyRejectController {
	return &SupplyRejectController{
		service: services.NewSupplyRejectService(),
	}
}

func (c *SupplyRejectController) CreateReject(ctx *gin.Context) {
	var req dtos.CreateSupplyRejectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	reject, err := c.service.CreateReject(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, reject)
}

func (c *SupplyRejectController) GetRejects(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetRejects(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Data": results, "Total": total})
}

func (c *SupplyRejectController) GetRejectsBySupply(ctx *gin.Context) {
	results, err := c.service.GetRejectsBySupply(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Data": results})
}

func (c *SupplyRejectController) GetReject(ctx *gin.Context) {
	result, err := c.service.GetReject(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Reject record not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *SupplyRejectController) UpdateReject(ctx *gin.Context) {
	var req dtos.UpdateSupplyRejectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateReject(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Reject record updated successfully"})
}

func (c *SupplyRejectController) DeleteReject(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteReject(ctx.Param("id"), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Reject record deleted successfully"})
}
