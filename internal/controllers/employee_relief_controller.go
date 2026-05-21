package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type EmployeeReliefController struct {
	service *services.EmployeeReliefService
}

func NewEmployeeReliefController() *EmployeeReliefController {
	return &EmployeeReliefController{
		service: services.NewEmployeeReliefService(),
	}
}

func (c *EmployeeReliefController) Create(ctx *gin.Context) {
	var req dtos.CreateEmployeeReliefRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	relief, err := c.service.Create(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, _ := c.service.Get(utils.Uint64ToString(relief.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *EmployeeReliefController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))
	employeeID := ctx.Query("employee_id")

	results, total, err := c.service.List(employeeID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *EmployeeReliefController) Get(ctx *gin.Context) {
	result, err := c.service.Get(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee relief record not found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *EmployeeReliefController) Update(ctx *gin.Context) {
	var req dtos.UpdateEmployeeReliefRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.Update(ctx.Param("id"), req, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Employee relief record updated successfully"})
}

func (c *EmployeeReliefController) Delete(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if err := c.service.Delete(ctx.Param("id"), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Employee relief record deleted successfully"})
}
