package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
)

type MobileLoanController struct {
	svc *services.MobileLoanService
}

func NewMobileLoanController() *MobileLoanController {
	return &MobileLoanController{svc: services.NewMobileLoanService()}
}

func (h *MobileLoanController) userID(c *gin.Context) uint64 {
	return c.MustGet("user_id").(uint64)
}

func respondServiceError(c *gin.Context, err error, data interface{}) {
	status := http.StatusInternalServerError
	var se *services.ServiceError
	if errors.As(err, &se) {
		status = se.Status
	}
	body := gin.H{"error": err.Error()}
	if data != nil {
		body["data"] = data
	}
	c.JSON(status, body)
}

func parseApplicationID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan request id"})
		return 0, false
	}
	return id, true
}

// RequestLoan serves POST /api/loan-module/mobile/loan-requests.
func (h *MobileLoanController) RequestLoan(c *gin.Context) {
	var req dtos.MobileLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.New().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationError(err)})
		return
	}
	resp, err := h.svc.RequestLoan(req, h.userID(c))
	if err != nil {
		respondServiceError(c, err, nil)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": resp, "message": "OTP sent to the member's registered phone"})
}

// ResendOTP serves POST /api/loan-module/mobile/loan-requests/:id/resend-otp.
func (h *MobileLoanController) ResendOTP(c *gin.Context) {
	id, ok := parseApplicationID(c)
	if !ok {
		return
	}
	resp, err := h.svc.ResendOTP(id, h.userID(c))
	if err != nil {
		respondServiceError(c, err, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp, "message": "A new OTP has been sent"})
}

// VerifyOTP serves POST /api/loan-module/mobile/loan-requests/:id/verify-otp.
func (h *MobileLoanController) VerifyOTP(c *gin.Context) {
	id, ok := parseApplicationID(c)
	if !ok {
		return
	}
	var req dtos.VerifyMobileLoanOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.New().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "otp must be exactly 6 digits"})
		return
	}
	resp, err := h.svc.VerifyOTP(id, req.OTP, h.userID(c))
	if err != nil {
		respondServiceError(c, err, resp)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// Quote serves GET /api/loan-module/mobile/quote?member_id=&amount= : the advance price (total repayable).
func (h *MobileLoanController) Quote(c *gin.Context) {
	memberID, _ := strconv.ParseUint(firstNonEmptyQuery(c, "member_id", "memberId"), 10, 64)
	amount, err := strconv.ParseFloat(c.Query("amount"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount is required"})
		return
	}
	q, err := h.svc.QuoteCashout(memberID, amount)
	if err != nil {
		respondServiceError(c, err, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": q})
}

// ListTransactions serves GET /api/loan-module/mobile/transactions?member_id=…&date_from=&date_to=
// The member's loan statement: cashouts sent, charges (processing, insurance, per-litre), interest,
// penalties and repayments, with totals.
func (h *MobileLoanController) ListTransactions(c *gin.Context) {
	raw := c.Query("member_id")
	if raw == "" {
		raw = c.Query("memberId")
	}
	memberID, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || memberID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "member_id is required"})
		return
	}
	lines, totals, err := services.NewLoanModuleService().ListMemberLoanTransactions(memberID, c.Query("date_from"), c.Query("date_to"))
	if err != nil {
		respondServiceError(c, err, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": lines, "total": len(lines), "totals": totals})
}

// ListRequests serves GET /api/loan-module/mobile/loan-requests?member_id=…
func (h *MobileLoanController) ListRequests(c *gin.Context) {
	raw := c.Query("member_id")
	if raw == "" {
		raw = c.Query("memberId")
	}
	memberID, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || memberID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "member_id is required"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	items, err := h.svc.ListRequests(memberID, limit)
	if err != nil {
		respondServiceError(c, err, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "total": len(items)})
}

// GetRequest serves GET /api/loan-module/mobile/loan-requests/:id.
func (h *MobileLoanController) GetRequest(c *gin.Context) {
	id, ok := parseApplicationID(c)
	if !ok {
		return
	}
	resp, err := h.svc.GetRequest(id)
	if err != nil {
		respondServiceError(c, err, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}
