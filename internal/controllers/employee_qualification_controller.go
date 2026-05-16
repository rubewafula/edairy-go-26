package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type EmployeeQualificationController struct {
	service *services.EmployeeQualificationService
}

func NewEmployeeQualificationController() *EmployeeQualificationController {
	return &EmployeeQualificationController{
		service: services.NewEmployeeQualificationService(),
	}
}

func (c *EmployeeQualificationController) Create(ctx *gin.Context) {
	var req dtos.CreateEmployeeQualificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extracting userID from Gin context (set by auth middleware)
	userID := ctx.GetUint64("userID")

	res, err := c.service.CreateQualification(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (c *EmployeeQualificationController) List(ctx *gin.Context) {
	employeeID := ctx.Query("employee_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	results, total, err := c.service.GetEmployeeQualifications(employeeID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  results,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (c *EmployeeQualificationController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.service.GetQualification(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Qualification not found"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *EmployeeQualificationController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateEmployeeQualificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint64("userID")

	if err := c.service.UpdateQualification(id, req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Qualification updated successfully"})
}

func (c *EmployeeQualificationController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetUint64("userID")

	if err := c.service.DeleteQualification(id, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Qualification deleted successfully"})
}
