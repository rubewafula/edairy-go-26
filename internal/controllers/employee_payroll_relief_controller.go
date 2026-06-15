package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"gorm.io/gorm"
)

type EmployeePayrollReliefController struct {
	service *services.EmployeePayrollReliefService
}

func NewEmployeePayrollReliefController() *EmployeePayrollReliefController {
	return &EmployeePayrollReliefController{
		service: services.NewEmployeePayrollReliefService(),
	}
}

func (c *EmployeePayrollReliefController) List(ctx *gin.Context) {
	employeeID := ctx.Query("employee_id")
	payrollID := ctx.Query("payroll_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	results, total, err := c.service.GetPayrollReliefs(employeeID, payrollID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  results,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (c *EmployeePayrollReliefController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.service.GetPayrollRelief(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Payroll relief not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *EmployeePayrollReliefController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetUint64("user_id")

	if err := c.service.DeletePayrollRelief(id, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Payroll relief deleted successfully"})
}
