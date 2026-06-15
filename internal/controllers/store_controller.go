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

type StoreController struct {
	service *services.StoreService
}

func NewStoreController() *StoreController {
	return &StoreController{
		service: services.NewStoreService(),
	}
}

func (c *StoreController) CreateStore(ctx *gin.Context) {
	var req dtos.CreateStoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[StoreController.CreateStore] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[StoreController.CreateStore] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	store, err := c.service.CreateStore(req)
	if err != nil {
		log.Printf("[StoreController.CreateStore] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create store"})
		return
	}
	ctx.JSON(http.StatusCreated, store)
}

func (c *StoreController) GetStores(ctx *gin.Context) {
	stores, total, err := c.service.GetStores()
	if err != nil {
		log.Printf("[StoreController.GetStores] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stores"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": stores, "total": total})
}

func (c *StoreController) GetStore(ctx *gin.Context) {
	store, err := c.service.GetStore(ctx.Param("id"))
	if err != nil {
		log.Printf("[StoreController.GetStore] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}
	ctx.JSON(http.StatusOK, store)
}

func (c *StoreController) UpdateStore(ctx *gin.Context) {
	var req dtos.UpdateStoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[StoreController.UpdateStore] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[StoreController.UpdateStore] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateStore(ctx.Param("id"), req); err != nil {
		log.Printf("[StoreController.UpdateStore] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update store"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Store updated successfully"})
}

func (c *StoreController) DeleteStore(ctx *gin.Context) {
	if err := c.service.DeleteStore(ctx.Param("id")); err != nil {
		log.Printf("[StoreController.DeleteStore] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete store"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Store deleted successfully"})
}

func (c *StoreController) ImportStoreStock(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Printf("[StoreController.ImportStoreStock] File Upload Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required for import"})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.ImportStoreStock(file, userID); err != nil {
		log.Printf("[StoreController.ImportStoreStock] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate store stock import"})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": "Store stock import started in the background."})
}

func (c *StoreController) GetImportErrors(ctx *gin.Context) {
	importIDStr := ctx.Param("importid")
	importID, _ := strconv.ParseUint(importIDStr, 10, 64)
	errors, err := c.service.GetImportErrors(importID)
	if err != nil {
		log.Printf("[StoreController.GetImportErrors] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve import errors"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": errors})
}
