package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type ActivityLogController struct {
	service *services.ActivityLogService
}

func NewActivityLogController() *ActivityLogController {
	return &ActivityLogController{
		service: services.NewActivityLogService(),
	}
}

func (ctrl *ActivityLogController) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	logs, total, err := ctrl.service.GetLogs(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (ctrl *ActivityLogController) GetLog(c *gin.Context) {
	id := c.Param("id")
	log, err := ctrl.service.GetLog(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}
