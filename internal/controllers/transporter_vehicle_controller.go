package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type TransporterVehicleController struct {
	service *services.TransporterVehicleService
}

func NewTransporterVehicleController() *TransporterVehicleController {
	return &TransporterVehicleController{
		service: services.NewTransporterVehicleService(),
	}
}

func (c *TransporterVehicleController) CreateVehicle(ctx *gin.Context) {
	var req dtos.CreateTransporterVehicleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	vehicle, err := c.service.CreateVehicle(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	response := dtos.TransporterVehicleResponse{
		ID:             vehicle.ID,
		TransporterID:  vehicle.TransporterID,
		RegistrationNo: vehicle.RegistrationNo,
		VehicleType:    vehicle.VehicleType,
		CapacityLitres: vehicle.CapacityLitres,
		Active:         vehicle.Active,
		CreatedAt:      vehicle.CreatedAt,
		UpdatedAt:      vehicle.UpdatedAt,
	}
	ctx.JSON(http.StatusCreated, response)
}

func (c *TransporterVehicleController) GetVehicles(ctx *gin.Context) {
	transporterID := ctx.Query("transporter_id")
	vehicles, total, err := c.service.GetVehicles(transporterID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": vehicles, "total": total})
}

func (c *TransporterVehicleController) GetVehicle(ctx *gin.Context) {
	vehicle, err := c.service.GetVehicle(ctx.Param("id")) // Now returns dtos.TransporterVehicleResponse
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Vehicle not found"})
		return
	}
	ctx.JSON(http.StatusOK, vehicle)
}

func (c *TransporterVehicleController) UpdateVehicle(ctx *gin.Context) {
	var req dtos.UpdateTransporterVehicleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateVehicle(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Vehicle updated successfully"})
}

func (c *TransporterVehicleController) DeleteVehicle(ctx *gin.Context) {
	if err := c.service.DeleteVehicle(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Vehicle deleted successfully"})
}
