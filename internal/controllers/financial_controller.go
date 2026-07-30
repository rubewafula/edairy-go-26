package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"github.com/rubewafula/edairy-go-26/internal/utils"
)

type FinancialController struct {
	posting *services.FinancialPostingService
	reports *services.FinancialReportService
	periods *services.FinancialPeriodService
	budget  *services.BudgetService
	bankRec *services.BankReconciliationService
}

func NewFinancialController() *FinancialController {
	return &FinancialController{
		posting: services.NewFinancialPostingService(),
		reports: services.NewFinancialReportService(),
		periods: services.NewFinancialPeriodService(),
		budget:  services.NewBudgetService(),
		bankRec: services.NewBankReconciliationService(),
	}
}

func (h *FinancialController) PostTransaction(c *gin.Context) {
	var req dtos.PostFinancialTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("user_id").(uint64)
	if req.IdempotencyKey == "" {
		req.IdempotencyKey = c.GetHeader("Idempotency-Key")
	}
	result, err := h.posting.PostFinancialTransaction(req, userID)
	if err != nil {
		log.Println("FinancialController.PostTransaction:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": result})
}

func (h *FinancialController) ReverseTransaction(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dtos.ReverseTransactionRequest
	_ = c.ShouldBindJSON(&req)
	userID := c.MustGet("user_id").(uint64)
	rev, err := h.posting.ReverseByTransactionID(id, userID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rev})
}

func (h *FinancialController) GetGeneralLedger(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	accountID, _ := strconv.ParseUint(c.DefaultQuery("account_id", "0"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	result, err := h.reports.GetGeneralLedger(from, to, accountID, page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *FinancialController) GetCashFlow(c *gin.Context) {
	result, err := h.reports.GetCashFlowStatement(c.Query("from"), c.Query("to"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *FinancialController) GetMemberStatement(c *gin.Context) {
	memberID, err := strconv.ParseUint(c.Param("member_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member_id"})
		return
	}
	result, err := h.reports.GetMemberStatement(memberID, c.Query("from"), c.Query("to"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *FinancialController) GetLoanStatement(c *gin.Context) {
	loanID, err := strconv.ParseUint(c.Param("loan_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan_id"})
		return
	}
	result, err := h.reports.GetLoanStatement(loanID, c.Query("from"), c.Query("to"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *FinancialController) ListPeriods(c *gin.Context) {
	periods, err := h.periods.ListPeriods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": periods})
}

func (h *FinancialController) CreatePeriod(c *gin.Context) {
	var req dtos.CreateFinancialPeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("user_id").(uint64)
	start := utils.ParseDate(req.StartDate)
	end := utils.ParseDate(req.EndDate)
	p, err := h.periods.CreatePeriod(req.Name, start, end, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": p})
}

func (h *FinancialController) ClosePeriod(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := c.MustGet("user_id").(uint64)
	if err := h.periods.ClosePeriod(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "period closed"})
}

func (h *FinancialController) GetBudgetVsActual(c *gin.Context) {
	result, _ := h.budget.GetBudgetVsActual(c.Query("from"), c.Query("to"))
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *FinancialController) GetBankReconciliation(c *gin.Context) {
	result, _ := h.bankRec.GetReconciliationStatus()
	c.JSON(http.StatusOK, gin.H{"data": result})
}
