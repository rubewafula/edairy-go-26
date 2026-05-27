package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"gorm.io/gorm"
)

type MemberPayDateRangeController struct {
	service *services.MemberPayDateRangeService
}

func NewMemberPayDateRangeController() *MemberPayDateRangeController {
	return &MemberPayDateRangeController{
		service: services.NewMemberPayDateRangeService(),
	}
}

func (c *MemberPayDateRangeController) Create(ctx *gin.Context) {
	var req dtos.CreateMemberPayDateRangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint64("user_id")

	res, err := c.service.Create(req, userID)
	if err != nil {
		log.Printf("[MemberPayDateRangeController.Create] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *MemberPayDateRangeController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	res, total, err := c.service.List(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *MemberPayDateRangeController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.service.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Member pay date range not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *MemberPayDateRangeController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateMemberPayDateRangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint64("user_id")

	if err := c.service.Update(id, req, userID); err != nil {
		log.Printf("[MemberPayDateRangeController.Update] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Member pay date range updated successfully"})
}

func (c *MemberPayDateRangeController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		log.Printf("[MemberPayDateRangeController.Delete] Error deleting range %s: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Member pay date range deleted successfully"})
}
