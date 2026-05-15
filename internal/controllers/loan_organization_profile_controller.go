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

type LoanOrganizationProfileController struct {
	service *services.LoanOrganizationProfileService
}

func NewLoanOrganizationProfileController() *LoanOrganizationProfileController {
	return &LoanOrganizationProfileController{
		service: services.NewLoanOrganizationProfileService(),
	}
}

func (c *LoanOrganizationProfileController) CreateLoanOrganizationProfile(ctx *gin.Context) {
	var req dtos.CreateLoanOrganizationProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	profile, err := c.service.CreateLoanOrganizationProfile(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response, _ := c.service.GetLoanOrganizationProfile(utils.Uint64ToString(profile.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *LoanOrganizationProfileController) GetLoanOrganizationProfiles(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetLoanOrganizationProfiles(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *LoanOrganizationProfileController) GetLoanOrganizationProfile(ctx *gin.Context) {
	result, err := c.service.GetLoanOrganizationProfile(ctx.Param("id"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Loan organization profile not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *LoanOrganizationProfileController) UpdateLoanOrganizationProfile(ctx *gin.Context) {
	var req dtos.UpdateLoanOrganizationProfileRequest
	id := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateLoanOrganizationProfile(id, req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Loan organization profile updated successfully"})
}

func (c *LoanOrganizationProfileController) DeleteLoanOrganizationProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteLoanOrganizationProfile(id, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Loan organization profile deleted successfully"})
}
