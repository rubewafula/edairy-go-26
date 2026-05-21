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

type LivestockController struct {
	service *services.LivestockService
}

func NewLivestockController() *LivestockController {
	return &LivestockController{service: services.NewLivestockService()}
}

func (c *LivestockController) CreateLivestocks(ctx *gin.Context) {
	var req dtos.CreateLivestockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}
	userID := ctx.GetUint64("user_id")
	res, err := c.service.CreateLivestock(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *LivestockController) GetLivestocks(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	res, total, err := c.service.GetLivestocks(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *LivestockController) GetLivestock(ctx *gin.Context) {
	res, err := c.service.GetLivestock(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Livestock not found"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *LivestockController) UpdateLivestocks(ctx *gin.Context) {
	var req dtos.UpdateLivestockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateLivestock(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Livestock updated successfully"})
}

func (c *LivestockController) DeleteLivestocks(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteLivestock(ctx.Param("id"), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Livestock deleted successfully"})
}

func (c *LivestockController) CreateCategory(ctx *gin.Context) {
	var req dtos.CreateLivestockCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}
	res, err := c.service.CreateCategory(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *LivestockController) GetCategories(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	res, total, err := c.service.GetCategories(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *LivestockController) GetCategory(ctx *gin.Context) {
	res, err := c.service.GetCategory(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *LivestockController) UpdateCategory(ctx *gin.Context) {
	var req dtos.UpdateLivestockCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateCategory(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

func (c *LivestockController) DeleteCategory(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteCategory(ctx.Param("id"), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// Breed CRUD
func (c *LivestockController) CreateBreed(ctx *gin.Context) {
	var req dtos.CreateLivestockBreedRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}
	res, err := c.service.CreateBreed(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *LivestockController) GetBreeds(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	res, total, err := c.service.GetBreeds(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *LivestockController) GetBreed(ctx *gin.Context) {
	res, err := c.service.GetBreed(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Breed not found"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *LivestockController) UpdateBreed(ctx *gin.Context) {
	var req dtos.UpdateLivestockBreedRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateBreed(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Breed updated successfully"})
}

func (c *LivestockController) DeleteBreed(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteBreed(ctx.Param("id"), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Breed deleted successfully"})
}

func (c *LivestockController) UpdateFeeding(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateLivestockFeedingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateFeeding(id, req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Livestock feeding record updated successfully"})
}

func (c *LivestockController) CreateFeeding(ctx *gin.Context) {
	var req dtos.CreateLivestockFeedingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	res, err := c.service.CreateFeeding(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *LivestockController) GetFeeding(ctx *gin.Context) {
	res, err := c.service.GetFeeding(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Feeding record not found"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *LivestockController) DeleteFeeding(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteFeeding(id, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Livestock feeding record deleted successfully"})
}

func (c *LivestockController) GetFeedings(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	id := ctx.Query("livestock_id")
	res, total, err := c.service.GetFeedings(id, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *LivestockController) CreateSale(ctx *gin.Context) {
	var req dtos.CreateLivestockSaleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}
	userID := ctx.GetUint64("user_id")
	res, err := c.service.CreateSale(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, _ := c.service.GetSale(utils.Uint64ToString(res.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *LivestockController) GetSales(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	id := ctx.Query("livestock_id")
	res, total, err := c.service.GetSales(id, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *LivestockController) GetSale(ctx *gin.Context) {
	res, err := c.service.GetSale(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Sale record not found"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *LivestockController) UpdateSale(ctx *gin.Context) {
	var req dtos.UpdateLivestockSaleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateSale(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Sale updated successfully"})
}

func (c *LivestockController) DeleteSale(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteSale(ctx.Param("id"), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Sale deleted successfully"})
}

func (c *LivestockController) CreateHealth(ctx *gin.Context) {
	var req dtos.CreateLivestockHealthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}
	userID := ctx.GetUint64("user_id")
	res, err := c.service.CreateHealth(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, _ := c.service.GetHealthRecord(utils.Uint64ToString(res.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *LivestockController) GetHealthRecords(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	id := ctx.Query("livestock_id")
	res, total, err := c.service.GetHealthRecords(id, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *LivestockController) GetHealthRecord(ctx *gin.Context) {
	res, err := c.service.GetHealthRecord(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Health record not found"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *LivestockController) UpdateHealthRecord(ctx *gin.Context) {
	var req dtos.UpdateLivestockHealthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateHealthRecord(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Health record updated successfully"})
}

func (c *LivestockController) DeleteHealthRecord(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteHealthRecord(ctx.Param("id"), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Health record deleted successfully"})
}

// Production CRUD
func (c *LivestockController) CreateProduction(ctx *gin.Context) {
	var req dtos.CreateLivestockProductionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}
	userID := ctx.GetUint64("user_id")
	res, err := c.service.CreateProduction(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *LivestockController) GetProductionRecords(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	id := ctx.Query("livestock_id")
	res, total, err := c.service.GetProductionRecords(id, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *LivestockController) GetProductionRecord(ctx *gin.Context) {
	res, err := c.service.GetProductionRecord(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Production record not found"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *LivestockController) UpdateProductionRecord(ctx *gin.Context) {
	var req dtos.UpdateLivestockProductionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateProductionRecord(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Production record updated successfully"})
}

func (c *LivestockController) DeleteProductionRecord(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteProductionRecord(ctx.Param("id"), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Production record deleted successfully"})
}

// Deaths CRUD
func (c *LivestockController) CreateDeath(ctx *gin.Context) {
	var req dtos.CreateLivestockDeathRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	res, err := c.service.CreateDeath(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *LivestockController) GetDeaths(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	id := ctx.Query("livestock_id")
	res, total, err := c.service.GetDeaths(id, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *LivestockController) GetDeath(ctx *gin.Context) {
	res, err := c.service.GetDeath(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Death record not found"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *LivestockController) UpdateDeath(ctx *gin.Context) {
	var req dtos.UpdateLivestockDeathRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateDeath(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Death record updated successfully"})
}

func (c *LivestockController) DeleteDeath(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteDeath(ctx.Param("id"), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Death record deleted successfully"})
}

// Movements CRUD
func (c *LivestockController) CreateMovement(ctx *gin.Context) {
	var req dtos.CreateLivestockMovementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	res, err := c.service.CreateMovement(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *LivestockController) GetMovements(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	id := ctx.Query("livestock_id")
	res, total, err := c.service.GetMovements(id, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *LivestockController) GetMovement(ctx *gin.Context) {
	res, err := c.service.GetMovement(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Movement record not found"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *LivestockController) UpdateMovement(ctx *gin.Context) {
	var req dtos.UpdateLivestockMovementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateMovement(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Movement updated successfully"})
}

func (c *LivestockController) DeleteMovement(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteMovement(ctx.Param("id"), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Movement deleted successfully"})
}

// Photos CRUD
func (c *LivestockController) CreatePhoto(ctx *gin.Context) {
	var req dtos.CreateLivestockPhotoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	res, err := c.service.CreatePhoto(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *LivestockController) GetPhotos(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	id := ctx.Query("livestock_id")
	res, total, err := c.service.GetPhotos(id, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *LivestockController) GetPhoto(ctx *gin.Context) {
	res, err := c.service.GetPhoto(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *LivestockController) UpdatePhoto(ctx *gin.Context) {
	var req dtos.UpdateLivestockPhotoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdatePhoto(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Photo updated successfully"})
}

func (c *LivestockController) DeletePhoto(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeletePhoto(ctx.Param("id"), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}

// Weights CRUD
func (c *LivestockController) CreateWeight(ctx *gin.Context) {
	var req dtos.CreateLivestockWeightRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	res, err := c.service.CreateWeight(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *LivestockController) GetWeightRecords(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	id := ctx.Query("livestock_id")
	res, total, err := c.service.GetWeightRecords(id, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "total": total})
}

func (c *LivestockController) GetWeightRecord(ctx *gin.Context) {
	res, err := c.service.GetWeightRecord(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Weight record not found"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *LivestockController) UpdateWeightRecord(ctx *gin.Context) {
	var req dtos.UpdateLivestockWeightRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateWeightRecord(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Weight record updated successfully"})
}

func (c *LivestockController) DeleteWeightRecord(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteWeightRecord(ctx.Param("id"), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Weight record deleted successfully"})
}
