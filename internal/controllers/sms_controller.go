package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	validator "github.com/rubewafula/edairy-go-26/internal/validators"
	"gorm.io/gorm"
)

type SMSController struct {
	service *services.SMSService
}

func NewSMSController() *SMSController {
	return &SMSController{service: services.NewSMSService()}
}

func (c *SMSController) CreateGroup(ctx *gin.Context) {
	var req dtos.CreateSMSGroupRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("[SMSController.CreateGroup] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID := ctx.GetUint64("user_id")

	group, err := c.service.CreateGroup(req, req.ContactsList, userID)
	if err != nil {
		log.Printf("[SMSController.CreateGroup] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SMS group"})
		return
	}
	ctx.JSON(http.StatusCreated, group)
}

func (c *SMSController) GetGroup(ctx *gin.Context) {
	id := ctx.Param("id")
	group, err := c.service.GetGroup(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS group not found"})
			return
		}
		log.Printf("[SMSController.GetGroup] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SMS group"})
		return
	}
	ctx.JSON(http.StatusOK, group)
}

func (c *SMSController) GetGroups(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))
	results, total, err := c.service.GetGroups(page, limit)
	if err != nil {
		log.Printf("[SMSController.GetGroups] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SMS groups"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *SMSController) UpdateGroup(ctx *gin.Context) {
	var req dtos.CreateSMSGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SMSController.UpdateGroup] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := c.service.UpdateGroup(ctx.Param("id"), req); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS group not found"})
			return
		}
		log.Printf("[SMSController.UpdateGroup] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update SMS group"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "SMS group updated successfully"})
}

func (c *SMSController) DeleteGroup(ctx *gin.Context) {
	if err := c.service.DeleteGroup(ctx.Param("id")); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS group not found"})
			return
		}
		log.Printf("[SMSController.DeleteGroup] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SMS group"})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (c *SMSController) CreateContact(ctx *gin.Context) {
	var req dtos.CreateSMSContactRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SMSController.CreateContact] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	contact, err := c.service.CreateContact(req)
	if err != nil {
		log.Printf("[SMSController.CreateContact] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SMS contact"})
		return
	}
	ctx.JSON(http.StatusCreated, contact)
}

func (c *SMSController) GetContacts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))
	results, total, err := c.service.GetContacts(page, limit)
	if err != nil {
		log.Printf("[SMSController.GetContacts] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SMS contacts"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *SMSController) GetContact(ctx *gin.Context) {
	id := ctx.Param("id")
	contact, err := c.service.GetContact(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS contact not found"})
			return
		}
		log.Printf("[SMSController.GetContact] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SMS contact"})
		return
	}
	ctx.JSON(http.StatusOK, contact)
}

func (c *SMSController) UpdateContact(ctx *gin.Context) {
	var req dtos.CreateSMSContactRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SMSController.UpdateContact] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	if err := c.service.UpdateContact(ctx.Param("id"), req); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS contact not found"})
			return
		}
		log.Printf("[SMSController.UpdateContact] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update SMS contact"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "SMS contact updated successfully"})
}

func (c *SMSController) DeleteContact(ctx *gin.Context) {
	if err := c.service.DeleteContact(ctx.Param("id")); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS contact not found"})
			return
		}
		log.Printf("[SMSController.DeleteContact] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SMS contact"})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (c *SMSController) GetContactsByGroup(ctx *gin.Context) {
	contacts, err := c.service.GetContactsByGroup(ctx.Param("id"))
	if err != nil {
		log.Printf("[SMSController.GetContactsByGroup] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve contacts for this group"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": contacts})
}

func (c *SMSController) SendMessage(ctx *gin.Context) {
	var req dtos.SendSMSRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SMSController.SendMessage] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	msg, err := c.service.SendSMS(req)
	if err != nil {
		log.Printf("[SMSController.SendMessage] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send SMS"})
		return
	}
	ctx.JSON(http.StatusOK, msg)
}

func (c *SMSController) GetQueue(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))
	results, total, err := c.service.GetQueue(page, limit)
	if err != nil {
		log.Printf("[SMSController.GetQueue] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SMS queue"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *SMSController) CreateProvider(ctx *gin.Context) {
	var req dtos.CreateSMSProviderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SMSController.CreateProvider] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[SMSController.CreateProvider] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	provider, err := c.service.CreateProvider(req, userID)
	if err != nil {
		log.Printf("[SMSController.CreateProvider] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SMS provider"})
		return
	}
	ctx.JSON(http.StatusCreated, provider)
}

func (c *SMSController) CreateTemplate(ctx *gin.Context) {
	var req dtos.CreateSMSTemplateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SMSController.CreateTemplate] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	userID := ctx.GetUint64("user_id")
	tpl, err := c.service.CreateTemplate(req, userID)
	if err != nil {
		log.Printf("[SMSController.CreateTemplate] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SMS template"})
		return
	}
	ctx.JSON(http.StatusCreated, tpl)
}

func (c *SMSController) GetTemplate(ctx *gin.Context) {
	result, err := c.service.GetTemplate(ctx.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS template not found"})
			return
		}
		log.Printf("[SMSController.GetTemplate] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SMS template"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *SMSController) UpdateTemplate(ctx *gin.Context) {
	var req dtos.CreateSMSTemplateRequest // Assuming CreateSMSTemplateRequest is reused for update
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SMSController.UpdateTemplate] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		log.Printf("[SMSController.UpdateTemplate] Validation Error: %v", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateTemplate(ctx.Param("id"), req, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS template not found"})
			return
		}
		log.Printf("[SMSController.UpdateTemplate] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update SMS template"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "SMS template updated successfully"})
}

func (c *SMSController) DeleteTemplate(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetUint64("user_id")

	if err := c.service.DeleteTemplate(id, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS template not found"})
			return
		}
		log.Printf("[SMSController.DeleteTemplate] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SMS template"})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// GetTemplate is missing in the provided code, adding it here.
// UpdateTemplate is missing in the provided code, adding it here.
// DeleteTemplate is missing in the provided code, adding it here.

func (c *SMSController) GetProviders(ctx *gin.Context) {
	results, err := c.service.GetProviders()
	if err != nil {
		log.Printf("[SMSController.GetProviders] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SMS providers"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results})
}

func (c *SMSController) GetProvider(ctx *gin.Context) {
	result, err := c.service.GetProvider(ctx.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS provider not found"})
			return
		}
		log.Printf("[SMSController.GetProvider] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SMS provider"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *SMSController) UpdateProvider(ctx *gin.Context) {
	var req dtos.CreateSMSProviderRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SMSController.UpdateProvider] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateProvider(ctx.Param("id"), req, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS provider not found"})
			return
		}
		log.Printf("[SMSController.UpdateProvider] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update SMS provider"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "SMS provider updated successfully"})
}

func (c *SMSController) DeleteProvider(ctx *gin.Context) {
	if err := c.service.DeleteProvider(ctx.Param("id")); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS provider not found"})
			return
		}
		log.Printf("[SMSController.DeleteProvider] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SMS provider"})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (c *SMSController) GetMessages(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))
	results, total, err := c.service.GetMessages(page, limit)
	if err != nil {
		log.Printf("[SMSController.GetMessages] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SMS messages"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *SMSController) GetMessage(ctx *gin.Context) {
	id := ctx.Param("id")
	msg, err := c.service.GetMessage(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS message not found"})
			return
		}
		log.Printf("[SMSController.GetMessage] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SMS message"})
		return
	}
	ctx.JSON(http.StatusOK, msg)
}

func (c *SMSController) CreateMessage(ctx *gin.Context) {
	var req dtos.SendSMSRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SMSController.CreateMessage] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	msg, err := c.service.CreateMessage(req)
	if err != nil {
		log.Printf("[SMSController.CreateMessage] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SMS message"})
		return
	}
	ctx.JSON(http.StatusCreated, msg)
}

func (c *SMSController) UpdateMessage(ctx *gin.Context) {
	var req dtos.SendSMSRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SMSController.UpdateMessage] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	if err := c.service.UpdateMessage(ctx.Param("id"), req); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS message not found"})
			return
		}
		log.Printf("[SMSController.UpdateMessage] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update SMS message"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "SMS message updated successfully"})
}

func (c *SMSController) DeleteMessage(ctx *gin.Context) {
	if err := c.service.DeleteMessage(ctx.Param("id")); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS message not found"})
			return
		}
		log.Printf("[SMSController.DeleteMessage] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SMS message"})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (c *SMSController) GetTemplates(ctx *gin.Context) {
	results, err := c.service.GetTemplates()
	if err != nil {
		log.Printf("[SMSController.GetTemplates] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SMS templates"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results})
}

func (c *SMSController) GetOutboxesByCampaign(ctx *gin.Context) {
	results, err := c.service.GetSMSOutboxesByCampaign(ctx.Param("id"))
	if err != nil {
		log.Printf("[SMSController.GetOutboxesByCampaign] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve campaign outbox"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results})
}

func (c *SMSController) CreateOutbox(ctx *gin.Context) {
	var req dtos.CreateSMSOutboxRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SMSController.CreateOutbox] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	userID := ctx.GetUint64("user_id")
	outbox, err := c.service.CreateSMSOutbox(req, userID)
	if err != nil {
		log.Printf("[SMSController.CreateOutbox] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SMS outbox entry"})
		return
	}
	ctx.JSON(http.StatusCreated, outbox)
}

func (c *SMSController) GetAllOutboxes(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))
	results, total, err := c.service.GetAllSMSOutboxes(page, limit)
	if err != nil {
		log.Printf("[SMSController.GetAllOutboxes] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SMS outboxes"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

func (c *SMSController) GetOutbox(ctx *gin.Context) {
	id := ctx.Param("id")
	outbox, err := c.service.GetSMSOutbox(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS outbox not found"})
			return
		}
		log.Printf("[SMSController.GetOutbox] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SMS outbox entry"})
		return
	}
	ctx.JSON(http.StatusOK, outbox)
}

func (c *SMSController) UpdateOutbox(ctx *gin.Context) {
	var req dtos.UpdateSMSOutboxRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SMSController.UpdateOutbox] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateSMSOutbox(ctx.Param("id"), req, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS outbox not found"})
			return
		}
		log.Printf("[SMSController.UpdateOutbox] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update SMS outbox"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "SMS outbox updated successfully"})
}

func (c *SMSController) DeleteOutbox(ctx *gin.Context) {
	if err := c.service.DeleteSMSOutbox(ctx.Param("id")); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "SMS outbox not found"})
			return
		}
		log.Printf("[SMSController.DeleteOutbox] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SMS outbox"})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// GetInAppConfigs handles the retrieval of all In-App SMS configurations with pagination.
func (c *SMSController) GetInAppConfigs(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("Page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("Limit", "10"))

	results, total, err := c.service.GetInAppConfigs(page, limit)
	if err != nil {
		log.Printf("[SMSController.GetInAppConfigs] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve in-app configurations"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  results,
		"total": total,
	})
}

// GetInAppConfig handles retrieving a single in-app configuration by ID.
func (c *SMSController) GetInAppConfig(ctx *gin.Context) {
	id := ctx.Param("id")
	config, err := c.service.GetInAppConfig(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found"})
			return
		}
		log.Printf("[SMSController.GetInAppConfig] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve configuration"})
		return
	}
	ctx.JSON(http.StatusOK, config)
}

// CreateInAppConfig handles the creation of a new in-app SMS configuration.
func (c *SMSController) CreateInAppConfig(ctx *gin.Context) {
	var req dtos.CreateSMSInAppConfigurationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SMSController.CreateInAppConfig] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID := ctx.GetUint64("user_id")
	config, err := c.service.CreateInAppConfig(req, userID)
	if err != nil {
		log.Printf("[SMSController.CreateInAppConfig] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create configuration"})
		return
	}
	ctx.JSON(http.StatusCreated, config)
}

// UpdateInAppConfig handles updating an existing in-app configuration.
func (c *SMSController) UpdateInAppConfig(ctx *gin.Context) {
	var req dtos.UpdateSMSInAppConfigurationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[SMSController.UpdateInAppConfig] Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID := ctx.GetUint64("user_id")
	if err := c.service.UpdateInAppConfig(ctx.Param("id"), req, userID); err != nil {
		log.Printf("[SMSController.UpdateInAppConfig] Service Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update configuration"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully"})
}

// DeleteInAppConfig handles soft-deleting an in-app configuration.
func (c *SMSController) DeleteInAppConfig(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if err := c.service.DeleteInAppConfig(ctx.Param("id"), userID); err != nil {
		log.Printf("[SMSController.DeleteInAppConfig] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete configuration"})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// GetImportErrors returns the validation/processing errors for a specific SMS group import session.
func (c *SMSController) GetImportErrors(ctx *gin.Context) {
	idStr := ctx.Param("id")
	importID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid import session ID"})
		return
	}

	importErrors, err := c.service.GetImportErrors(importID)
	if err != nil {
		log.Printf("[SMSController.GetImportErrors] Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve import errors"})
		return
	}

	ctx.JSON(http.StatusOK, importErrors)
}
