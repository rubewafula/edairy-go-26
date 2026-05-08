package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type LoanController struct {
	service *services.LoanService
}

func NewLoanController() *LoanController {
	return &LoanController{
		service: services.NewLoanService(),
	}
}

func (c *LoanController) CreateLoan(ctx *gin.Context) {
	var req dtos.CreateLoanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	loan, err := c.service.CreateLoan(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, loan)
}

func (c *LoanController) GetLoans(ctx *gin.Context) {
	loans, total, err := c.service.GetLoans()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": loans, "total": total})
}

func (c *LoanController) GetLoan(ctx *gin.Context) {
	loan, err := c.service.GetLoan(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}
	ctx.JSON(http.StatusOK, loan)
}

func (c *LoanController) UpdateLoan(ctx *gin.Context) {
	var req dtos.UpdateLoanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateLoan(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Loan updated successfully"})
}

func (c *LoanController) DeleteLoan(ctx *gin.Context) {
	if err := c.service.DeleteLoan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Loan deleted successfully"})
}
