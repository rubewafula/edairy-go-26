package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type CustomerClassController struct {
	service *services.CustomerClassService
}

func NewCustomerClassController() *CustomerClassController {
	return &CustomerClassController{
		service: services.NewCustomerClassService(),
	}
}

func (c *CustomerClassController) CreateClass(ctx *gin.Context) {
	var req dtos.CreateCustomerClassRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	class, err := c.service.CreateCustomerClass(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, class)
}

func (c *CustomerClassController) GetClasses(ctx *gin.Context) {
	classes, total, err := c.service.GetCustomerClasses()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": classes, "total": total})
}

func (c *CustomerClassController) GetClass(ctx *gin.Context) {
	class, err := c.service.GetCustomerClass(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}
	ctx.JSON(http.StatusOK, class)
}

func (c *CustomerClassController) UpdateClass(ctx *gin.Context) {
	var req dtos.UpdateCustomerClassRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateCustomerClass(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Class updated successfully"})
}

func (c *CustomerClassController) DeleteClass(ctx *gin.Context) {
	if err := c.service.DeleteCustomerClass(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Class deleted successfully"})
}
