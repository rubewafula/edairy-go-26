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

type CustomerTypeController struct {
	service *services.CustomerTypeService
}

func NewCustomerTypeController() *CustomerTypeController {
	return &CustomerTypeController{
		service: services.NewCustomerTypeService(),
	}
}

func (c *CustomerTypeController) CreateType(ctx *gin.Context) {
	var req dtos.CreateCustomerTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	customerType, err := c.service.CreateCustomerType(req)
	if err != nil {
		log.Printf("[CustomerTypeController.CreateType] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, customerType)
}

func (c *CustomerTypeController) GetTypes(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	customerTypes, total, err := c.service.GetCustomerTypes(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": customerTypes, "total": total})
}

func (c *CustomerTypeController) GetType(ctx *gin.Context) {
	customerType, err := c.service.GetCustomerType(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Type not found"})
		return
	}
	ctx.JSON(http.StatusOK, customerType)
}

func (c *CustomerTypeController) UpdateType(ctx *gin.Context) {
	var req dtos.UpdateCustomerTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateCustomerType(ctx.Param("id"), req); err != nil {
		log.Printf("[CustomerTypeController.UpdateType] Error updating type %s: %v", ctx.Param("id"), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Type updated successfully"})
}

func (c *CustomerTypeController) DeleteType(ctx *gin.Context) {
	if err := c.service.DeleteCustomerType(ctx.Param("id")); err != nil {
		log.Printf("[CustomerTypeController.DeleteType] Error deleting type %s: %v", ctx.Param("id"), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Type deleted successfully"})
}
