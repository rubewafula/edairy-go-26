package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type MemberCreditController struct {
	svc *services.MemberCreditService
}

func NewMemberCreditController() *MemberCreditController {
	return &MemberCreditController{svc: services.NewMemberCreditService()}
}

func (h *MemberCreditController) userID(c *gin.Context) uint64 {
	return c.MustGet("user_id").(uint64)
}

func parseMemberID(raw string) (uint64, bool) {
	id, err := strconv.ParseUint(raw, 10, 64)
	return id, err == nil && id > 0
}

// GetMemberCreditLimit serves GET /api/member-credit-limit?member={id}, the path used by the mobile app.
func (h *MemberCreditController) GetMemberCreditLimit(c *gin.Context) {
	raw := c.Query("member")
	if raw == "" {
		raw = c.Query("member_id")
	}
	memberID, ok := parseMemberID(raw)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "member query parameter is required"})
		return
	}
	h.respondCreditLimit(c, memberID)
}

// GetCreditLimit serves GET /api/members/:id/credit-limit.
func (h *MemberCreditController) GetCreditLimit(c *gin.Context) {
	memberID, ok := parseMemberID(c.Param("id"))
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member id"})
		return
	}
	h.respondCreditLimit(c, memberID)
}

func (h *MemberCreditController) respondCreditLimit(c *gin.Context, memberID uint64) {
	resp, err := h.svc.GetCreditLimit(memberID, h.userID(c))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "member not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// ListCreditLimitLogs serves GET /api/members/:id/credit-limit-logs.
func (h *MemberCreditController) ListCreditLimitLogs(c *gin.Context) {
	memberID, ok := parseMemberID(c.Param("id"))
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member id"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	logs, total, err := h.svc.ListLogs(memberID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": logs, "total": total})
}

// RecalculateCreditLimit serves POST /api/members/:id/credit-limit/recalculate.
func (h *MemberCreditController) RecalculateCreditLimit(c *gin.Context) {
	memberID, ok := parseMemberID(c.Param("id"))
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member id"})
		return
	}
	h.respondCreditLimit(c, memberID)
}

// Reconcile serves POST /api/member-credit-limits/reconcile[?member_id=ID]: corrects stored credit limits
// so each unpaid journal's milk is credited exactly once if confirmed (0 if not), with a log per journal.
func (h *MemberCreditController) Reconcile(c *gin.Context) {
	memberID, _ := strconv.ParseUint(c.Query("member_id"), 10, 64)
	res, err := h.svc.ReconcileConfirmedMilkCredit(memberID, h.userID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// RecalculateAll serves POST /api/member-credit-limits/recalculate-all.
func (h *MemberCreditController) RecalculateAll(c *gin.Context) {
	queued, err := h.svc.RecalculateAll(h.userID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"data": gin.H{"members_queued": queued}})
}
