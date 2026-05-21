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

type MemberBankAccountController struct {
	service *services.MemberBankAccountService
}

func NewMemberBankAccountController() *MemberBankAccountController {
	return &MemberBankAccountController{
		service: services.NewMemberBankAccountService(),
	}
}

func (c *MemberBankAccountController) CreateAccount(ctx *gin.Context) {
	var req dtos.CreateMemberBankAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	account, err := c.service.CreateAccount(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response, _ := c.service.GetAccount(utils.Uint64ToString(account.ID))
	ctx.JSON(http.StatusCreated, response)
}

func (c *MemberBankAccountController) GetAccounts(ctx *gin.Context) {

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	memberNo := ctx.Query("member_no")
	firstName := ctx.Query("first_name")
	lastName := ctx.Query("last_name")
	bankName := ctx.Query("bank_name")
	accountNo := ctx.Query("account_number")

	accounts, total, err := c.service.GetAccounts(
		page,
		limit,
		memberNo,
		firstName,
		lastName,
		bankName,
		accountNo,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": accounts, "total": total})
}

func (c *MemberBankAccountController) GetAccount(ctx *gin.Context) {
	account, err := c.service.GetAccount(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	ctx.JSON(http.StatusOK, account)
}

func (c *MemberBankAccountController) UpdateAccount(ctx *gin.Context) {
	var req dtos.UpdateMemberBankAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	if err := c.service.UpdateAccount(ctx.Param("id"), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Account updated successfully"})
}

func (c *MemberBankAccountController) DeleteAccount(ctx *gin.Context) {
	if err := c.service.DeleteAccount(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
