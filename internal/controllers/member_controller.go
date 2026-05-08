package controllers

import (
	"net/http"

	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"

	"github.com/gin-gonic/gin"
)

type MemberController struct {
	service *services.MemberService
}

func NewMemberController() *MemberController {
	return &MemberController{
		service: services.NewMemberService(),
	}
}

// POST /users
func (c *MemberController) CreateMember(ctx *gin.Context) {

	var req dtos.CreateMemberRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(422, gin.H{
			"error": utils.FormatValidationError(err),
		})
		return
	}

	idFront, _ := ctx.FormFile("id_front_photo")
	idBack, _ := ctx.FormFile("id_back_photo")
	passport, _ := ctx.FormFile("passport_photo")

	member, err := c.service.CreateMember(
		ctx.Request.Context(),
		req,
		idFront,
		idBack,
		passport,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, member)
}

// GET /users
func (c *MemberController) GetMembers(ctx *gin.Context) {
	members, total, err := c.service.GetMembers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": members, "total": total})
}

// GET /users/:id
func (c *MemberController) GetMember(ctx *gin.Context) {
	member, err := c.service.GetMember(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, member)
}

func (c *MemberController) UpdateMember(ctx *gin.Context) {
	var req dtos.UpdateMemberRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.JSON(422, gin.H{
			"error": utils.FormatValidationError(err),
		})
		return
	}

	idFront, _ := ctx.FormFile("id_front_photo")
	idBack, _ := ctx.FormFile("id_back_photo")
	passport, _ := ctx.FormFile("passport_photo")

	err := c.service.UpdateMember(
		ctx.Param("id"),
		req,
		idFront,
		idBack,
		passport,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

func (c *MemberController) DeleteMember(ctx *gin.Context) {
	err := c.service.DeleteMember(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}
