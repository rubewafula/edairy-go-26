package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type LocationController struct {
	service *services.LocationService
}

func NewLocationController() *LocationController {
	return &LocationController{
		service: services.NewLocationService(),
	}
}

// Administrative Location APIs
func (c *LocationController) CreateLocation(ctx *gin.Context) {
	var req dtos.CreateLocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint64("user_id")
	location, err := c.service.CreateLocation(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, location)
}

func (c *LocationController) GetLocations(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Support 'size' parameter from client
	if size := ctx.Query("size"); size != "" {
		if s, err := strconv.Atoi(size); err == nil {
			limit = s
		}
	}

	results, total, err := c.service.GetLocations(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *LocationController) GetLocation(ctx *gin.Context) {
	result, err := c.service.GetLocation(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *LocationController) UpdateLocation(ctx *gin.Context) {
	var req dtos.UpdateLocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateLocation(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Location updated successfully"})
}

func (c *LocationController) DeleteLocation(ctx *gin.Context) {
	if err := c.service.DeleteLocation(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Location deleted successfully"})
}

// Counties
func (c *LocationController) GetCounties(ctx *gin.Context) {
	counties, err := c.service.GetCounties()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": counties})
}

func (c *LocationController) CreateCounty(ctx *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Code string `json:"code" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	county, err := c.service.CreateCounty(req.Name, req.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, county)
}

// SubCounties
func (c *LocationController) GetSubCounties(ctx *gin.Context) {
	countyID := ctx.Query("county_id")
	subCounties, err := c.service.GetSubCounties(countyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": subCounties})
}

func (c *LocationController) CreateSubCounty(ctx *gin.Context) {
	var req struct {
		CountyID uint64 `json:"county_id" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	subCounty, err := c.service.CreateSubCounty(req.CountyID, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, subCounty)
}

// Wards
func (c *LocationController) GetWards(ctx *gin.Context) {
	subCountyID := ctx.Query("sub_county_id")
	wards, err := c.service.GetWards(subCountyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": wards})
}

func (c *LocationController) CreateWard(ctx *gin.Context) {
	var req struct {
		SubCountyID uint64 `json:"sub_county_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ward, err := c.service.CreateWard(req.SubCountyID, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, ward)
}
