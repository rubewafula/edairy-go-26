package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
