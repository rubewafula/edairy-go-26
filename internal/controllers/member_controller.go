package controllers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
	"gorm.io/gorm"

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

	log.Printf("Received Content-Type: %s", ctx.ContentType())
	var req dtos.CreateMemberRequest

	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("Create Member Binding error: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {

		log.Printf("Create member, Validation error: %s", err.Error())
		ctx.JSON(422, gin.H{
			"error": utils.FormatValidationError(err),
		})
		return
	}
	userID := ctx.GetUint64("user_id")
	member, err := c.service.CreateMember(ctx.Request.Context(), req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, member)
}

// GET /users
// GET /members
func (c *MemberController) GetMembers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	memberNo := ctx.Query("member_no")
	primaryPhone := ctx.Query("primary_phone")
	memberTypeID := ctx.Query("member_type_id")
	routeID := ctx.Query("route_id")

	members, total, err := c.service.GetMembers(
		page,
		limit,
		memberNo,
		primaryPhone,
		memberTypeID,
		routeID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": members, "total": total})
}

// GET /members/:id
func (c *MemberController) GetMember(ctx *gin.Context) {
	member, err := c.service.GetMember(ctx.Param("id"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, member)
}

func (c *MemberController) UpdateMember(ctx *gin.Context) {
	var req dtos.UpdateMemberRequest

	if err := ctx.ShouldBind(&req); err != nil {
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

	userID := ctx.GetUint64("user_id")
	err := c.service.UpdateMember(
		ctx.Param("id"),
		req,
		userID,
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

func (c *MemberController) SuspendMember(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.SuspendMember(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Member suspended successfully"})
}

func (c *MemberController) ImportMembers(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.ImportMembers(file, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": "Member import started in the background. Check logs for status."})
}

func (c *MemberController) GetMemberImportErrors(ctx *gin.Context) {
	importIDStr := ctx.Param("importid")
	importID, _ := strconv.ParseUint(importIDStr, 10, 64)

	errors, err := c.service.GetImportErrors(importID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch import errors"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": errors})
}

func (c *MemberController) ExportMembers(ctx *gin.Context) {
	memberNo := ctx.Query("member_no")
	primaryPhone := ctx.Query("primary_phone")
	memberTypeID := ctx.Query("member_type_id")
	routeID := ctx.Query("route_id")
	gender := ctx.Query("gender")
	status := ctx.Query("status")
	reportType := ctx.DefaultQuery("format", "csv")

	userID := ctx.GetUint64("user_id")
	if err := c.service.ExportMembers(userID, memberNo, primaryPhone, memberTypeID, routeID, gender, status, reportType); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"message": "Member export started in the background. You will receive a notification when it's ready."})
}

func (c *MemberController) DownloadExportFile(ctx *gin.Context) {
	filename := filepath.Base(ctx.Param("filename"))
	filePath := filepath.Join("./storage/exports", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Export file not found"})
		return
	}

	ctx.File(filePath)
}
