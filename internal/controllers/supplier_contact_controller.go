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

type SupplierContactController struct {
	service *services.SupplierContactService
}

func NewSupplierContactController() *SupplierContactController {
	return &SupplierContactController{
		service: services.NewSupplierContactService(),
	}
}

func (c *SupplierContactController) CreateContact(ctx *gin.Context) {
	var req dtos.CreateSupplierContactRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SupplierContactController.CreateContact] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[SupplierContactController.CreateContact] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	contact, err := c.service.CreateContact(req, userID)
	if err != nil {
		log.Printf("[SupplierContactController.CreateContact] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create supplier contact"})
		return
	}
	response, _ := c.service.GetContact(utils.Uint64ToString(contact.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *SupplierContactController) GetContacts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetContacts(page, limit)
	if err != nil {
		log.Printf("[SupplierContactController.GetContacts] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supplier contacts"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *SupplierContactController) GetContact(ctx *gin.Context) {
	result, err := c.service.GetContact(ctx.Param("id"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Supplier contact not found"})
			return
		}
		log.Printf("[SupplierContactController.GetContact] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supplier contact"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *SupplierContactController) UpdateContact(ctx *gin.Context) {
	var req dtos.UpdateSupplierContactRequest
	id := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SupplierContactController.UpdateContact] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[SupplierContactController.UpdateContact] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateContact(id, req, userID); err != nil {
		log.Printf("[SupplierContactController.UpdateContact] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update supplier contact"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Supplier contact updated successfully"})
}

func (c *SupplierContactController) DeleteContact(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteContact(id, userID); err != nil {
		log.Printf("[SupplierContactController.DeleteContact] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete supplier contact"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Supplier contact deleted successfully"})
}

func (c *SupplierContactController) GetContactsBySupplier(ctx *gin.Context) {
	results, err := c.service.GetContactsBySupplier(ctx.Param("id"))
	if err != nil {
		log.Printf("[SupplierContactController.GetContactsBySupplier] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve contacts by supplier"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results})
}
