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

type TransporterDriverController struct {
	service *services.TransporterDriverService
}

func NewTransporterDriverController() *TransporterDriverController {
	return &TransporterDriverController{
		service: services.NewTransporterDriverService(),
	}
}

func (c *TransporterDriverController) CreateDriver(ctx *gin.Context) {
	var req dtos.CreateTransporterDriverRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[TransporterDriverController.CreateDriver] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[TransporterDriverController.CreateDriver] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	driver, err := c.service.CreateDriver(req)
	if err != nil {
		log.Printf("[TransporterDriverController.CreateDriver] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transporter driver"})
		return
	}

	response, _ := c.service.GetDriver(utils.Uint64ToString(driver.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *TransporterDriverController) GetDrivers(ctx *gin.Context) {
	drivers, total, err := c.service.GetDrivers() // Now returns dtos.TransporterDriverResponse
	if err != nil {
		log.Printf("[TransporterDriverController.GetDrivers] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transporter drivers"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": drivers, "total": total})
}

func (c *TransporterDriverController) GetDriver(ctx *gin.Context) {
	driver, err := c.service.GetDriver(ctx.Param("id"))
	if err != nil {
		log.Printf("[TransporterDriverController.GetDriver] Service Error: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Transporter driver not found"})
		return
	}
	ctx.JSON(http.StatusOK, driver)
}

func (c *TransporterDriverController) UpdateDriver(ctx *gin.Context) {
	var req dtos.UpdateTransporterDriverRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[TransporterDriverController.UpdateDriver] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[TransporterDriverController.UpdateDriver] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateDriver(ctx.Param("id"), req); err != nil {
		log.Printf("[TransporterDriverController.UpdateDriver] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transporter driver"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Driver updated successfully"})
}

func (c *TransporterDriverController) DeleteDriver(ctx *gin.Context) {
	if err := c.service.DeleteDriver(ctx.Param("id")); err != nil {
		log.Printf("[TransporterDriverController.DeleteDriver] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transporter driver"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Driver deleted successfully"})
}
