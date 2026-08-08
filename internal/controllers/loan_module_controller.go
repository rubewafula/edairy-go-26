package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
)

type LoanModuleController struct {
	svc *services.LoanModuleService
}

func NewLoanModuleController() *LoanModuleController {
	return &LoanModuleController{svc: services.NewLoanModuleService()}
}

func (h *LoanModuleController) parsePage(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 200 {
		limit = 50
	}
	return page, limit
}

func (h *LoanModuleController) userID(c *gin.Context) uint64 {
	return c.MustGet("user_id").(uint64)
}

// Products

func (h *LoanModuleController) CreateProduct(c *gin.Context) {
	var req dtos.CreateLoanProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.New().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationError(err)})
		return
	}
	p, err := h.svc.CreateProduct(req, h.userID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": p})
}

func (h *LoanModuleController) ListProducts(c *gin.Context) {
	page, limit := h.parsePage(c)
	active := c.Query("active") == "true"
	items, total, err := h.svc.ListProducts(page, limit, active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "total": total})
}

func (h *LoanModuleController) GetProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *LoanModuleController) UpdateProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dtos.UpdateLoanProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p, err := h.svc.UpdateProduct(id, req, h.userID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// Applications

func (h *LoanModuleController) CreateApplication(c *gin.Context) {
	var req dtos.CreateLoanApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	app, err := h.svc.CreateApplication(req, h.userID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": app})
}

func (h *LoanModuleController) ListApplications(c *gin.Context) {
	page, limit := h.parsePage(c)
	memberID, _ := strconv.ParseUint(c.Query("member_id"), 10, 64)
	items, total, err := h.svc.ListApplications(page, limit, c.Query("status"), memberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "total": total})
}

func (h *LoanModuleController) GetApplication(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	app, err := h.svc.GetApplication(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": app})
}

func (h *LoanModuleController) SubmitApplication(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	app, err := h.svc.SubmitApplication(id, h.userID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": app})
}

func (h *LoanModuleController) ApproveApplication(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dtos.ApproveLoanApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	app, err := h.svc.ApproveApplication(id, req, h.userID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": app})
}

func (h *LoanModuleController) RejectApplication(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dtos.RejectLoanApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	app, err := h.svc.RejectApplication(id, req, h.userID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": app})
}

// Contracts

func (h *LoanModuleController) ListContracts(c *gin.Context) {
	page, limit := h.parsePage(c)
	memberID, _ := strconv.ParseUint(c.Query("member_id"), 10, 64)
	items, total, err := h.svc.ListContracts(page, limit, c.Query("status"), memberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "total": total})
}

func (h *LoanModuleController) GetContract(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	contract, err := h.svc.GetContract(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "contract not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": contract})
}

func (h *LoanModuleController) ListDisbursementChannels(c *gin.Context) {
	all := c.Query("all") == "true" || c.Query("include_inactive") == "true"
	items, err := h.svc.ListDisbursementChannels(all)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "total": len(items)})
}

func (h *LoanModuleController) GetDisbursementChannel(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ch, err := h.svc.GetDisbursementChannel(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ch})
}

func (h *LoanModuleController) CreateDisbursementChannel(c *gin.Context) {
	var req dtos.CreateDisbursementChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.New().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationError(err)})
		return
	}
	ch, err := h.svc.CreateDisbursementChannel(req, h.userID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": ch})
}

func (h *LoanModuleController) UpdateDisbursementChannel(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dtos.UpdateDisbursementChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ch, err := h.svc.UpdateDisbursementChannel(id, req, h.userID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ch})
}

func (h *LoanModuleController) DeleteDisbursementChannel(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteDisbursementChannel(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"id": id, "deleted": true}})
}

func (h *LoanModuleController) GetDisbursementChannelConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	items, err := h.svc.GetDisbursementChannelConfig(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *LoanModuleController) UpdateDisbursementChannelConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dtos.UpdateDisbursementChannelConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdateDisbursementChannelConfig(id, req.Items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items, err := h.svc.GetDisbursementChannelConfig(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": req.Items, "updated": true})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "updated": true})
}

func (h *LoanModuleController) ConfirmDisbursement(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dtos.ConfirmLoanDisbursementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d, err := h.svc.ConfirmDisbursement(id, req, h.userID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": d})
}

func (h *LoanModuleController) WithdrawalCallback(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ProcessWithdrawalCallback(data); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "received", "note": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "processed"})
}

func (h *LoanModuleController) MpesaB2CResultCallback(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ProcessMpesaB2CResult(data); err != nil {
		c.JSON(http.StatusOK, gin.H{"ResultCode": 0, "ResultDesc": "Accepted", "note": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ResultCode": 0, "ResultDesc": "Success"})
}

func (h *LoanModuleController) MpesaB2CTimeoutCallback(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ProcessMpesaB2CTimeout(data); err != nil {
		c.JSON(http.StatusOK, gin.H{"ResultCode": 0, "ResultDesc": "Accepted", "note": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ResultCode": 0, "ResultDesc": "Success"})
}

func (h *LoanModuleController) AirtelDisbursementCallback(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ProcessAirtelDisbursementCallback(data); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "received", "note": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "processed"})
}

func (h *LoanModuleController) JengaMobileCallback(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ProcessJengaMobileCallback(data); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "received", "note": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "processed"})
}

func (h *LoanModuleController) GetDisbursementProviderStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	result, err := h.svc.GetDisbursementProviderStatus(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *LoanModuleController) GetDisbursementQuote(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	quote, err := h.svc.GetDisbursementQuote(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": quote})
}

func (h *LoanModuleController) DisburseContract(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dtos.DisburseLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.IdempotencyKey == "" {
		req.IdempotencyKey = c.GetHeader("Idempotency-Key")
	}
	d, err := h.svc.DisburseContract(id, req, h.userID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": d})
}

func (h *LoanModuleController) GetSchedule(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	contract, err := h.svc.GetContract(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "contract not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": contract.Installments})
}

func (h *LoanModuleController) RecordRepayment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dtos.RecordLoanRepaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.IdempotencyKey == "" {
		req.IdempotencyKey = c.GetHeader("Idempotency-Key")
	}
	rep, err := h.svc.RecordRepayment(id, req, h.userID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": rep})
}

func (h *LoanModuleController) GetStatement(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	stmt, err := h.svc.GetStatement(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stmt})
}

func (h *LoanModuleController) SettleContract(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dtos.SettleLoanRequest
	_ = c.ShouldBindJSON(&req)
	if c.Query("quote") == "true" {
		req.QuoteOnly = true
	}
	quote, err := h.svc.SettleContract(id, req, h.userID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": quote})
}

func (h *LoanModuleController) RestructureContract(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dtos.RestructureLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r, err := h.svc.RestructureContract(id, req, h.userID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}

func (h *LoanModuleController) WriteOffContract(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dtos.WriteOffLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r, err := h.svc.WriteOffContract(id, req, h.userID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}

func (h *LoanModuleController) AccrueInterest(c *gin.Context) {
	var req dtos.AccrueInterestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	count, err := h.svc.AccrueInterest(req, h.userID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"contracts_processed": count}})
}

func (h *LoanModuleController) RunPenalties(c *gin.Context) {
	count, err := h.svc.RunPenalties(utils.Now(), h.userID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"penalties_applied": count}})
}

// Reports

func (h *LoanModuleController) PortfolioReport(c *gin.Context) {
	r, err := h.svc.PortfolioReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}

func (h *LoanModuleController) AgingReport(c *gin.Context) {
	r, err := h.svc.AgingReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}

func (h *LoanModuleController) PARReport(c *gin.Context) {
	r, err := h.svc.PARReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}

func (h *LoanModuleController) RegisterReport(c *gin.Context) {
	status := c.Query("status")
	memberID, _ := strconv.ParseUint(c.Query("member_id"), 10, 64)
	r, err := h.svc.RegisterReport(status, memberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}

func (h *LoanModuleController) ApplicationsPipelineReport(c *gin.Context) {
	r, err := h.svc.ApplicationsPipelineReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}

func (h *LoanModuleController) DisbursementsReport(c *gin.Context) {
	status := c.Query("status")
	channelCode := c.Query("channel_code")
	r, err := h.svc.DisbursementsReport(status, channelCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}
