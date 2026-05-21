package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
)

type BankBranchController struct {
	service *services.BankBranchService
}

func NewBankBranchController() *BankBranchController {
	return &BankBranchController{
		service: services.NewBankBranchService(),
	}
}

func (c *BankBranchController) CreateBankBranch(ctx *gin.Context) {
	var req dtos.CreateBankBranchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	branch, err := c.service.CreateBankBranch(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response, _ := c.service.GetBankBranch(utils.Uint64ToString(branch.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *BankBranchController) GetBankBranches(ctx *gin.Context) {
	branches, total, err := c.service.GetBankBranches()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": branches, "total": total})
}

func (c *BankBranchController) GetBankBranch(ctx *gin.Context) {
	branch, err := c.service.GetBankBranch(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Bank branch not found"})
		return
	}
	ctx.JSON(http.StatusOK, branch)
}

func (c *BankBranchController) UpdateBankBranch(ctx *gin.Context) {
	var req dtos.UpdateBankBranchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateBankBranch(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Bank branch updated successfully"})
}

func (c *BankBranchController) DeleteBankBranch(ctx *gin.Context) {
	if err := c.service.DeleteBankBranch(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Bank branch deleted successfully"})
}
