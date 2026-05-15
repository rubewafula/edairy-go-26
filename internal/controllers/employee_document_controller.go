package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
	"gorm.io/gorm"
)

type EmployeeDocumentController struct {
	service *services.EmployeeDocumentService
}

func NewEmployeeDocumentController() *EmployeeDocumentController {
	return &EmployeeDocumentController{
		service: services.NewEmployeeDocumentService(),
	}
}

func (c *EmployeeDocumentController) CreateEmployeeDocument(ctx *gin.Context) {
	var req dtos.CreateEmployeeDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}
	userID := ctx.GetUint64("user_id")
	document, err := c.service.CreateEmployeeDocument(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, document)
}

func (c *EmployeeDocumentController) GetEmployeeDocuments(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))
	results, total, err := c.service.GetEmployeeDocuments(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *EmployeeDocumentController) GetEmployeeDocument(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.GetEmployeeDocument(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee document not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *EmployeeDocumentController) UpdateEmployeeDocument(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateEmployeeDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateEmployeeDocument(id, req, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee document not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Employee document updated successfully"})
}

func (c *EmployeeDocumentController) DeleteEmployeeDocument(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteEmployeeDocument(id, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee document not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Employee document deleted successfully"})
}
