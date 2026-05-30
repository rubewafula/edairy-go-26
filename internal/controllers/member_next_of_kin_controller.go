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

type MemberNextOfKinController struct {
	service *services.MemberNextOfKinService
}

func NewMemberNextOfKinController() *MemberNextOfKinController {
	return &MemberNextOfKinController{
		service: services.NewMemberNextOfKinService(),
	}
}

// CreateMemberNextOfKin godoc
// @Summary Create a new member next of kin
// @Description Create a new member next of kin record
// @Tags Member Next Of Kin
// @Accept json
// @Produce json
// @Param nextOfKin body dtos.CreateMemberNextOfKinRequest true "Member Next Of Kin object to be created"
// @Success 201 {object} models.MemberNextOfKin
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /member-next-of-kins [post]
func (c *MemberNextOfKinController) CreateMemberNextOfKin(ctx *gin.Context) {
	var req dtos.CreateMemberNextOfKinRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.MustGet("user_id").(uint64)

	nextOfKin, err := c.service.CreateMemberNextOfKin(req, userID)
	if err != nil {
		log.Printf("[MemberNextOfKinController.CreateMemberNextOfKin] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, nextOfKin)
}

// GetMemberNextOfKins godoc
// @Summary Get all member next of kins
// @Description Retrieve a list of all member next of kins, with optional filtering by member ID
// @Tags Member Next Of Kin
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param member_id query string false "Filter by Member ID"
// @Success 200 {array} dtos.MemberNextOfKinResponse
// @Failure 500 {object} map[string]string
// @Router /member-next-of-kins [get]
func (c *MemberNextOfKinController) GetMemberNextOfKins(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	memberID := ctx.Query("member_id")

	nextOfKins, total, err := c.service.GetMemberNextOfKins(page, limit, memberID)
	if err != nil {
		log.Printf("[MemberNextOfKinController.GetMemberNextOfKins] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": nextOfKins, "total": total})
}

// GetMemberNextOfKin godoc
// @Summary Get a single member next of kin by ID
// @Description Retrieve a member next of kin record by its ID
// @Tags Member Next Of Kin
// @Produce json
// @Param id path string true "Member Next Of Kin ID"
// @Success 200 {object} dtos.MemberNextOfKinResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /member-next-of-kins/{id} [get]
func (c *MemberNextOfKinController) GetMemberNextOfKin(ctx *gin.Context) {
	id := ctx.Param("id")
	nextOfKin, err := c.service.GetMemberNextOfKin(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Member next of kin not found"})
			return
		}
		log.Printf("[MemberNextOfKinController.GetMemberNextOfKin] Error getting next of kin %s: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, nextOfKin)
}

// UpdateMemberNextOfKin godoc
// @Summary Update an existing member next of kin
// @Description Update a member next of kin record with the provided ID and details
// @Tags Member Next Of Kin
// @Accept json
// @Produce json
// @Param id path string true "Member Next Of Kin ID"
// @Param nextOfKin body dtos.UpdateMemberNextOfKinRequest true "Updated Member Next Of Kin object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /member-next-of-kins/{id} [put]
func (c *MemberNextOfKinController) UpdateMemberNextOfKin(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateMemberNextOfKinRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.MustGet("user_id").(uint64)

	err := c.service.UpdateMemberNextOfKin(id, req, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Member next of kin not found"})
			return
		}
		log.Printf("[MemberNextOfKinController.UpdateMemberNextOfKin] Error updating next of kin %s: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Member next of kin updated successfully"})
}

// DeleteMemberNextOfKin godoc
// @Summary Delete a member next of kin
// @Description Delete a member next of kin record by its ID
// @Tags Member Next Of Kin
// @Produce json
// @Param id path string true "Member Next Of Kin ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /member-next-of-kins/{id} [delete]
func (c *MemberNextOfKinController) DeleteMemberNextOfKin(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.MustGet("user_id").(uint64)
	err := c.service.DeleteMemberNextOfKin(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Member next of kin not found"})
			return
		}
		log.Printf("[MemberNextOfKinController.DeleteMemberNextOfKin] Error deleting next of kin %s: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Member next of kin deleted successfully"})
}
