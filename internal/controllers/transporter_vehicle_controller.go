package controllers

import (
	"net/http"

	"log"

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
		log.Printf("[TransporterVehicleController.CreateVehicle] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[TransporterVehicleController.CreateVehicle] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	vehicle, err := c.service.CreateVehicle(req)
	if err != nil {
		log.Printf("[TransporterVehicleController.CreateVehicle] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transporter vehicle"})
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
		log.Printf("[TransporterVehicleController.GetVehicles] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transporter vehicles"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": vehicles, "total": total})
}

func (c *TransporterVehicleController) GetVehicle(ctx *gin.Context) {
	vehicle, err := c.service.GetVehicle(ctx.Param("id")) // Now returns dtos.TransporterVehicleResponse
	if err != nil {
		log.Printf("[TransporterVehicleController.GetVehicle] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Transporter vehicle not found"})
		return
	}
	ctx.JSON(http.StatusOK, vehicle)
}

func (c *TransporterVehicleController) UpdateVehicle(ctx *gin.Context) {
	var req dtos.UpdateTransporterVehicleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[TransporterVehicleController.UpdateVehicle] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[TransporterVehicleController.UpdateVehicle] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateVehicle(ctx.Param("id"), req); err != nil {
		log.Printf("[TransporterVehicleController.UpdateVehicle] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transporter vehicle"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Vehicle updated successfully"})
}

func (c *TransporterVehicleController) DeleteVehicle(ctx *gin.Context) {
	if err := c.service.DeleteVehicle(ctx.Param("id")); err != nil {
		log.Printf("[TransporterVehicleController.DeleteVehicle] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transporter vehicle"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Vehicle deleted successfully"})
}
