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

type SupplierCategoryController struct {
	service *services.SupplierCategoryService
}

func NewSupplierCategoryController() *SupplierCategoryController {
	return &SupplierCategoryController{
		service: services.NewSupplierCategoryService(),
	}
}

func (c *SupplierCategoryController) CreateCategory(ctx *gin.Context) {
	var req dtos.CreateSupplierCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SupplierCategoryController.CreateCategory] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[SupplierCategoryController.CreateCategory] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	category, err := c.service.CreateCategory(req, userID)
	if err != nil {
		log.Printf("[SupplierCategoryController.CreateCategory] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create supplier category"})
		return
	}
	ctx.JSON(http.StatusCreated, category)
}

func (c *SupplierCategoryController) GetCategories(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetCategories(page, limit)
	if err != nil {
		log.Printf("[SupplierCategoryController.GetCategories] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supplier categories"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *SupplierCategoryController) GetCategory(ctx *gin.Context) {
	result, err := c.service.GetCategory(ctx.Param("id"))
	if err != nil {
		log.Printf("[SupplierCategoryController.GetCategory] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Supplier category not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *SupplierCategoryController) UpdateCategory(ctx *gin.Context) {
	var req dtos.CreateSupplierCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SupplierCategoryController.UpdateCategory] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateCategory(ctx.Param("id"), req, userID); err != nil {
		log.Printf("[SupplierCategoryController.UpdateCategory] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update supplier category"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

func (c *SupplierCategoryController) DeleteCategory(ctx *gin.Context) {
	if err := c.service.DeleteCategory(ctx.Param("id")); err != nil {
		log.Printf("[SupplierCategoryController.DeleteCategory] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete supplier category"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
