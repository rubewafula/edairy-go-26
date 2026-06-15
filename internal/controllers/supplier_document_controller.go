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
	"gorm.io/gorm"
)

type SupplierDocumentController struct {
	service *services.SupplierDocumentService
}

func NewSupplierDocumentController() *SupplierDocumentController {
	return &SupplierDocumentController{
		service: services.NewSupplierDocumentService(),
	}
}

func (c *SupplierDocumentController) CreateDocument(ctx *gin.Context) {
	var req dtos.CreateSupplierDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SupplierDocumentController.CreateDocument] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[SupplierDocumentController.CreateDocument] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	doc, err := c.service.CreateDocument(req, userID)
	if err != nil {
		log.Printf("[SupplierDocumentController.CreateDocument] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create supplier document"})
		return
	}
	response, _ := c.service.GetDocument(utils.Uint64ToString(doc.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *SupplierDocumentController) GetDocuments(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetDocuments(page, limit)
	if err != nil {
		log.Printf("[SupplierDocumentController.GetDocuments] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supplier documents"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *SupplierDocumentController) GetDocument(ctx *gin.Context) {
	result, err := c.service.GetDocument(ctx.Param("id"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Supplier document not found"})
			return
		}
		log.Printf("[SupplierDocumentController.GetDocument] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supplier document"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *SupplierDocumentController) UpdateDocument(ctx *gin.Context) {
	var req dtos.UpdateSupplierDocumentRequest
	id := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SupplierDocumentController.UpdateDocument] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[SupplierDocumentController.UpdateDocument] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateDocument(id, req, userID); err != nil {
		log.Printf("[SupplierDocumentController.UpdateDocument] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update supplier document"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Supplier document updated successfully"})
}

func (c *SupplierDocumentController) DeleteDocument(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteDocument(id, userID); err != nil {
		log.Printf("[SupplierDocumentController.DeleteDocument] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete supplier document"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Supplier document deleted successfully"})
}

func (c *SupplierDocumentController) GetDocumentsBySupplier(ctx *gin.Context) {
	results, err := c.service.GetDocumentsBySupplier(ctx.Param("id"))
	if err != nil {
		log.Printf("[SupplierDocumentController.GetDocumentsBySupplier] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents by supplier"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results})
}

func (c *SupplierDocumentController) VerifyDocument(ctx *gin.Context) {
	var req dtos.VerifySupplierDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SupplierDocumentController.VerifyDocument] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.VerifyDocument(ctx.Param("id"), req, userID); err != nil {
		log.Printf("[SupplierDocumentController.VerifyDocument] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document verification status"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Document verification status updated"})
}
