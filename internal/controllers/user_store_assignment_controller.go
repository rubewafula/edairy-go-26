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

type UserStoreAssignmentController struct {
	service *services.UserStoreAssignmentService
}

func NewUserStoreAssignmentController() *UserStoreAssignmentController {
	return &UserStoreAssignmentController{
		service: services.NewUserStoreAssignmentService(),
	}
}

func (c *UserStoreAssignmentController) CreateAssignment(ctx *gin.Context) {
	var req dtos.CreateUserStoreAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id") // Assuming user_id is set by auth middleware
	assignment, err := c.service.CreateUserStoreAssignment(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Fetch the full response DTO to include joined names
	response, err := c.service.GetUserStoreAssignment(utils.Uint64ToString(assignment.ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, response)
}

func (c *UserStoreAssignmentController) GetAssignments(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	userID := ctx.Query("user_id")
	storeID := ctx.Query("store_id")

	assignments, total, err := c.service.GetUserStoreAssignments(page, limit, userID, storeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": assignments, "total": total})
}

func (c *UserStoreAssignmentController) GetAssignment(ctx *gin.Context) {
	assignment, err := c.service.GetUserStoreAssignment(ctx.Param("id"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User store assignment not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, assignment)
}

func (c *UserStoreAssignmentController) UpdateAssignment(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateUserStoreAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id") // Assuming user_id is set by auth middleware
	if err := c.service.UpdateUserStoreAssignment(id, req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User store assignment updated successfully"})
}

func (c *UserStoreAssignmentController) DeleteAssignment(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetUint64("user_id") // Assuming user_id is set by auth middleware
	if err := c.service.DeleteUserStoreAssignment(id, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User store assignment not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User store assignment deleted successfully"})
}
