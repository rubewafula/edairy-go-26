package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type MemberPayrollDeductionService struct {
	notificationService *UINotificationService
}

func NewMemberPayrollDeductionService() *MemberPayrollDeductionService {
	return &MemberPayrollDeductionService{
		notificationService: NewUINotificationService(),
	}
}

// MemberPayrollDeductionFilters holds shared list/export filter values.
type MemberPayrollDeductionFilters struct {
	MemberID        string
	MemberNo        string
	PayrollID       string
	DeductionTypeID string
	DeductionMonth  string
	FiscalYear      string
	Settled         string
	Confirmed       string
}

type memberPayrollDeductionExportTitle struct {
	DairyName    string
	PlatformName string
	ReportTitle  string
}

func (s *MemberPayrollDeductionService) baseQuery() *gorm.DB {
	return db.DB.Table("member_payroll_deductions mpd").
		Select(`
			mpd.id, mpd.member_id, mr.member_no, CONCAT(mr.first_name, ' ', mr.last_name) AS member_name,
			mpd.deduction_month, mpd.fiscal_year, mpd.deduction_type_id, dt.description AS deduction_type_name,
			mpd.amount, mpd.priority, mpd.settled, mpd.transaction_date, mpd.date_captured,
			mpd.confirmed, mpd.payroll_id, mp.fiscal_period, mp.status AS payroll_status,
			mpd.reference, mpd.settlement_type,
			mpd.created_at, mpd.updated_at
		`).
		Joins("LEFT JOIN member_registrations mr ON mpd.member_id = mr.id").
		Joins("LEFT JOIN deduction_types dt ON mpd.deduction_type_id = dt.id").
		Joins("LEFT JOIN member_payrolls mp ON mpd.payroll_id = mp.id").
		Where("mpd.deleted_at IS NULL")
}

func (s *MemberPayrollDeductionService) enrichDefaultLoanFlags(rows []dtos.MemberPayrollDeductionResponse) {
	contractNos := make([]string, 0)
	seen := map[string]struct{}{}
	for _, row := range rows {
		ref := strings.TrimSpace(row.Reference)
		if !strings.HasPrefix(strings.ToUpper(ref), "LOAN-") {
			continue
		}
		contractNo := strings.TrimPrefix(ref, "LOAN-")
		contractNo = strings.TrimPrefix(contractNo, "loan-")
		if contractNo == "" {
			continue
		}
		if _, ok := seen[contractNo]; ok {
			continue
		}
		seen[contractNo] = struct{}{}
		contractNos = append(contractNos, contractNo)
	}
	if len(contractNos) == 0 {
		return
	}

	type loanFlag struct {
		ContractNo          string `gorm:"column:contract_no"`
		IsDefaultMemberLoan bool   `gorm:"column:is_default_member_loan"`
	}
	var flags []loanFlag
	_ = db.DB.Table("loan_contracts lc").
		Select("lc.contract_no, COALESCE(lp.is_default_member_loan, 0) AS is_default_member_loan").
		Joins("LEFT JOIN loan_products lp ON lp.id = lc.loan_product_id").
		Where("lc.contract_no IN ? AND lc.deleted_at IS NULL", contractNos).
		Scan(&flags).Error

	byContract := map[string]bool{}
	for _, f := range flags {
		byContract[f.ContractNo] = f.IsDefaultMemberLoan
	}

	for i := range rows {
		ref := strings.TrimSpace(rows[i].Reference)
		if !strings.HasPrefix(strings.ToUpper(ref), "LOAN-") {
			continue
		}
		contractNo := strings.TrimPrefix(ref, "LOAN-")
		contractNo = strings.TrimPrefix(contractNo, "loan-")
		if byContract[contractNo] {
			rows[i].IsDefaultMemberLoan = true
			rows[i].CannotReduceAmount = true
		}
	}
}

func (s *MemberPayrollDeductionService) isDefaultMemberLoanReference(reference string) bool {
	ref := strings.TrimSpace(reference)
	if !strings.HasPrefix(strings.ToUpper(ref), "LOAN-") {
		return false
	}
	contractNo := strings.TrimPrefix(ref, "LOAN-")
	contractNo = strings.TrimPrefix(contractNo, "loan-")
	var isDefault bool
	_ = db.DB.Table("loan_contracts lc").
		Select("COALESCE(lp.is_default_member_loan, 0)").
		Joins("LEFT JOIN loan_products lp ON lp.id = lc.loan_product_id").
		Where("lc.contract_no = ? AND lc.deleted_at IS NULL", contractNo).
		Limit(1).
		Scan(&isDefault).Error
	return isDefault
}

func (s *MemberPayrollDeductionService) applyFilters(query *gorm.DB, filters MemberPayrollDeductionFilters) *gorm.DB {
	if filters.MemberID != "" {
		query = query.Where("mpd.member_id = ?", filters.MemberID)
	}
	if filters.MemberNo != "" {
		query = query.Where("mr.member_no LIKE ?", "%"+filters.MemberNo+"%")
	}
	if filters.PayrollID != "" {
		query = query.Where("mpd.payroll_id = ?", filters.PayrollID)
	}
	if filters.DeductionTypeID != "" {
		query = query.Where("mpd.deduction_type_id = ?", filters.DeductionTypeID)
	}
	if filters.DeductionMonth != "" {
		query = query.Where("mpd.deduction_month = ?", filters.DeductionMonth)
	}
	if filters.FiscalYear != "" {
		query = query.Where("mpd.fiscal_year = ?", filters.FiscalYear)
	}
	if filters.Settled != "" {
		query = query.Where("mpd.settled = ?", filters.Settled)
	}
	if filters.Confirmed != "" {
		query = query.Where("mpd.confirmed = ?", filters.Confirmed)
	}
	return query
}

// GetMemberPayrollDeductions retrieves a list of member payroll deductions with optional filtering and pagination.
func (s *MemberPayrollDeductionService) GetMemberPayrollDeductions(page, limit int, filters MemberPayrollDeductionFilters) ([]dtos.MemberPayrollDeductionResponse, int64, error) {
	var deductions []dtos.MemberPayrollDeductionResponse
	var total int64

	query := s.applyFilters(s.baseQuery(), filters)
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Order("mpd.id DESC").Scan(&deductions).Error
	if err != nil {
		return deductions, total, err
	}
	s.enrichDefaultLoanFlags(deductions)
	return deductions, total, nil
}

// GetMemberPayrollDeduction retrieves a single member payroll deduction by its ID.
func (s *MemberPayrollDeductionService) GetMemberPayrollDeduction(id string) (*dtos.MemberPayrollDeductionResponse, error) {
	var deduction dtos.MemberPayrollDeductionResponse
	query := s.baseQuery().
		Where("mpd.id = ?", id).
		Limit(1)

	err := query.Scan(&deduction).Error
	if err != nil {
		return nil, err
	}
	if deduction.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	rows := []dtos.MemberPayrollDeductionResponse{deduction}
	s.enrichDefaultLoanFlags(rows)
	return &rows[0], nil
}

// UpdateMemberPayrollDeductionAmount reviews/updates a deduction amount while payroll is still editable.
// Default member loan deductions cannot be reduced.
func (s *MemberPayrollDeductionService) UpdateMemberPayrollDeductionAmount(id string, req dtos.UpdateMemberPayrollDeductionRequest, userID uint64) (*dtos.MemberPayrollDeductionResponse, error) {
	newAmount, err := strconv.ParseFloat(strings.TrimSpace(req.Amount), 64)
	if err != nil || newAmount < 0 {
		return nil, fmt.Errorf("amount must be a valid non-negative number")
	}

	var deduction models.MemberPayrollDeduction
	if err := db.DB.Where("id = ? AND deleted_at IS NULL", id).First(&deduction).Error; err != nil {
		return nil, err
	}

	var payroll models.MemberPayroll
	if err := db.DB.Where("id = ? AND deleted_at IS NULL", deduction.PayrollID).First(&payroll).Error; err != nil {
		return nil, fmt.Errorf("parent payroll not found")
	}

	status := strings.ToLower(strings.TrimSpace(payroll.Status))
	editable := map[string]bool{"draft": true, "incomplete": true, "rejected": true}
	if !editable[status] {
		return nil, fmt.Errorf("payroll deductions can only be reviewed while payroll status is draft, incomplete, or rejected (current: %s)", payroll.Status)
	}

	oldAmount, _ := strconv.ParseFloat(strings.TrimSpace(deduction.Amount), 64)
	isDefaultLoan := s.isDefaultMemberLoanReference(deduction.Reference)
	if isDefaultLoan && newAmount+0.0001 < oldAmount {
		return nil, fmt.Errorf("default member loan deductions cannot be reduced during payroll review")
	}

	amountStr := fmt.Sprintf("%.2f", newAmount)
	if err := db.DB.Model(&deduction).Updates(map[string]interface{}{
		"amount":     amountStr,
		"updated_by": userID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return nil, err
	}

	// Keep payslip totals in sync for this member.
	var total float64
	_ = db.DB.Table("member_payroll_deductions").
		Select("COALESCE(SUM(CAST(amount AS DECIMAL(15,2))), 0)").
		Where("payroll_id = ? AND member_id = ? AND deleted_at IS NULL", deduction.PayrollID, deduction.MemberID).
		Scan(&total).Error

	var payslip models.MemberPayslip
	if err := db.DB.Where("payroll_id = ? AND member_id = ? AND deleted_at IS NULL", deduction.PayrollID, deduction.MemberID).First(&payslip).Error; err == nil {
		net := payslip.GrossPay - total
		_ = db.DB.Model(&payslip).Updates(map[string]interface{}{
			"total_deductions": total,
			"net_pay":          net,
			"updated_by":       userID,
			"updated_at":       time.Now(),
		}).Error
	}

	var payrollTotal float64
	_ = db.DB.Table("member_payroll_deductions").
		Select("COALESCE(SUM(CAST(amount AS DECIMAL(15,2))), 0)").
		Where("payroll_id = ? AND deleted_at IS NULL", deduction.PayrollID).
		Scan(&payrollTotal).Error
	var payrollGross float64
	_ = db.DB.Table("member_payslips").
		Select("COALESCE(SUM(gross_pay), 0)").
		Where("payroll_id = ? AND deleted_at IS NULL", deduction.PayrollID).
		Scan(&payrollGross).Error
	_ = db.DB.Model(&payroll).Updates(map[string]interface{}{
		"total_deductions": payrollTotal,
		"net_pay":          payrollGross - payrollTotal,
		"updated_by":       userID,
		"updated_at":       time.Now(),
	}).Error

	return s.GetMemberPayrollDeduction(id)
}

// ExportDeductions initiates a background process to export member payroll deductions.
func (s *MemberPayrollDeductionService) ExportDeductions(userID uint64, filters MemberPayrollDeductionFilters, reportType string) error {
	go s.processDeductionExportInBackground(userID, filters, reportType)
	return nil
}

func (s *MemberPayrollDeductionService) processDeductionExportInBackground(userID uint64, filters MemberPayrollDeductionFilters, reportType string) {
	var results []dtos.MemberPayrollDeductionResponse

	query := s.applyFilters(s.baseQuery(), filters)
	if err := query.Order("mpd.id DESC").Scan(&results).Error; err != nil {
		log.Printf("[MemberPayrollDeductionService.processDeductionExportInBackground] Query error: %v", err)
		return
	}

	title := s.buildExportTitle(filters, results)

	var fileData []byte
	var err error
	ext := "csv"
	if strings.ToLower(reportType) == "pdf" {
		ext = "pdf"
		fileData, err = s.generateDeductionPDF(results, title)
	} else {
		fileData, err = s.generateDeductionCSV(results, title)
	}

	if err != nil {
		log.Printf("[MemberPayrollDeductionService.processDeductionExportInBackground] Generation error: %v", err)
		return
	}

	exportDir := "./storage/exports"
	os.MkdirAll(exportDir, 0755)
	filename := fmt.Sprintf("payroll_deductions_export_%d.%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(exportDir, filename)

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		log.Printf("[MemberPayrollDeductionService.processDeductionExportInBackground] File error: %v", err)
		return
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Deductions Export Complete",
		Message:          fmt.Sprintf("Your payroll deductions export (%s) is ready for download.", strings.ToUpper(ext)),
		NotificationType: "SUCCESS",
		DownloadLink:     fmt.Sprintf("/api/member-payroll-deductions/export/download/%s", filename),
	})
}

func (s *MemberPayrollDeductionService) resolveMemberLabel(filters MemberPayrollDeductionFilters, results []dtos.MemberPayrollDeductionResponse) string {
	if filters.MemberID == "" && strings.TrimSpace(filters.MemberNo) == "" {
		return ""
	}

	for _, row := range results {
		memberNo := strings.TrimSpace(row.MemberNo)
		names := strings.TrimSpace(row.MemberName)
		if memberNo == "" && names == "" {
			continue
		}
		if memberNo != "" && names != "" {
			return memberNo + " " + names
		}
		if memberNo != "" {
			return memberNo
		}
		return names
	}

	var member struct {
		MemberNo  string `gorm:"column:member_no"`
		FirstName string `gorm:"column:first_name"`
		LastName  string `gorm:"column:last_name"`
	}
	query := db.DB.Table("member_registrations").Select("member_no, first_name, last_name").Where("deleted_at IS NULL")
	if filters.MemberID != "" {
		query = query.Where("id = ?", filters.MemberID)
	} else {
		query = query.Where("member_no = ?", filters.MemberNo)
	}
	if err := query.Limit(1).Scan(&member).Error; err != nil {
		return strings.TrimSpace(filters.MemberNo)
	}

	names := strings.TrimSpace(strings.TrimSpace(member.FirstName) + " " + strings.TrimSpace(member.LastName))
	memberNo := strings.TrimSpace(member.MemberNo)
	if memberNo != "" && names != "" {
		return memberNo + " " + names
	}
	if memberNo != "" {
		return memberNo
	}
	if names != "" {
		return names
	}
	return strings.TrimSpace(filters.MemberNo)
}

func (s *MemberPayrollDeductionService) resolveDeductionTypeName(filters MemberPayrollDeductionFilters, results []dtos.MemberPayrollDeductionResponse) string {
	if filters.DeductionTypeID == "" {
		return ""
	}
	for _, row := range results {
		if name := strings.TrimSpace(row.DeductionTypeName); name != "" {
			return name
		}
	}
	var name string
	_ = db.DB.Table("deduction_types").
		Select("description").
		Where("id = ? AND deleted_at IS NULL", filters.DeductionTypeID).
		Limit(1).
		Scan(&name).Error
	return strings.TrimSpace(name)
}

func (s *MemberPayrollDeductionService) resolvePayrollLabel(filters MemberPayrollDeductionFilters, results []dtos.MemberPayrollDeductionResponse) string {
	if filters.PayrollID == "" {
		return ""
	}
	for _, row := range results {
		if period := strings.TrimSpace(row.FiscalPeriod); period != "" {
			return period
		}
	}
	var period string
	_ = db.DB.Table("member_payrolls").
		Select("fiscal_period").
		Where("id = ? AND deleted_at IS NULL", filters.PayrollID).
		Limit(1).
		Scan(&period).Error
	return strings.TrimSpace(period)
}

func (s *MemberPayrollDeductionService) buildExportTitle(filters MemberPayrollDeductionFilters, results []dtos.MemberPayrollDeductionResponse) memberPayrollDeductionExportTitle {
	var org struct {
		RegisteredName string `gorm:"column:registered_name"`
		TradingName    string `gorm:"column:trading_name"`
	}
	db.DB.Table("organization_details").Select("registered_name, trading_name").First(&org)

	dairyName := strings.TrimSpace(org.TradingName)
	if dairyName == "" {
		dairyName = strings.TrimSpace(org.RegisteredName)
	}
	if dairyName == "" {
		dairyName = "Dev Dairy"
	}

	reportTitle := "Member Payroll Deductions"
	if payrollLabel := s.resolvePayrollLabel(filters, results); payrollLabel != "" {
		reportTitle = fmt.Sprintf("Member Payroll Deductions - Payroll %s", payrollLabel)
	}
	if memberLabel := s.resolveMemberLabel(filters, results); memberLabel != "" {
		reportTitle = fmt.Sprintf("%s - Member %s", reportTitle, memberLabel)
	}
	if typeName := s.resolveDeductionTypeName(filters, results); typeName != "" {
		reportTitle = fmt.Sprintf("%s - %s", reportTitle, typeName)
	}

	return memberPayrollDeductionExportTitle{
		DairyName:    dairyName,
		PlatformName: "Edairy Africa Ltd",
		ReportTitle:  reportTitle,
	}
}

func (s *MemberPayrollDeductionService) generateDeductionCSV(results []dtos.MemberPayrollDeductionResponse, title memberPayrollDeductionExportTitle) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	_ = writer.Write([]string{title.DairyName})
	_ = writer.Write([]string{title.PlatformName})
	_ = writer.Write([]string{title.ReportTitle})
	_ = writer.Write([]string{})

	headers := []string{"Member No", "Member Name", "Month", "Year", "Type", "Amount", "Payroll Period", "Reference", "Settled", "Transaction Date"}
	_ = writer.Write(headers)

	var totalAmount float64

	for _, d := range results {
		amt, _ := strconv.ParseFloat(d.Amount, 64)
		totalAmount += amt

		transDate := ""
		if d.TransactionDate != nil {
			transDate = d.TransactionDate.Format("2006-01-02")
		}

		_ = writer.Write([]string{
			d.MemberNo,
			d.MemberName,
			d.DeductionMonth,
			fmt.Sprintf("%d", d.FiscalYear),
			d.DeductionTypeName,
			d.Amount,
			d.FiscalPeriod,
			d.Reference,
			d.Settled,
			transDate,
		})
	}

	_ = writer.Write([]string{
		"TOTALS", "", "", "", "",
		fmt.Sprintf("%.2f", totalAmount),
		"", "", "", "",
	})
	_ = writer.Write([]string{})
	_ = writer.Write([]string{"powered by edairy africa"})

	writer.Flush()
	return buf.Bytes(), writer.Error()
}

func (s *MemberPayrollDeductionService) generateDeductionPDF(results []dtos.MemberPayrollDeductionResponse, title memberPayrollDeductionExportTitle) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetFooterFunc(func() {
		pdf.SetY(-12)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 8, "powered by edairy africa", "", 0, "C", false, 0, "")
	})
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 9, title.DairyName, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 7, title.PlatformName, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 11)
	pdf.MultiCell(0, 6, title.ReportTitle, "", "C", false)
	pdf.Ln(4)

	pdf.SetFont("Arial", "B", 8)
	headers := []string{"M-No", "Member Name", "Month", "Year", "Type", "Amount", "Payroll Period", "Reference", "Settled", "Trans Date"}
	widths := []float64{22, 45, 18, 14, 35, 22, 28, 40, 16, 22}

	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	var totalAmount float64

	pdf.SetFont("Arial", "", 8)
	for _, d := range results {
		amt, _ := strconv.ParseFloat(d.Amount, 64)
		totalAmount += amt

		transDate := ""
		if d.TransactionDate != nil {
			transDate = d.TransactionDate.Format("2006-01-02")
		}

		pdf.CellFormat(widths[0], 8, d.MemberNo, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 8, d.MemberName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 8, d.DeductionMonth, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 8, fmt.Sprintf("%d", d.FiscalYear), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[4], 8, d.DeductionTypeName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[5], 8, d.Amount, "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[6], 8, d.FiscalPeriod, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[7], 8, d.Reference, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[8], 8, d.Settled, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[9], 8, transDate, "1", 0, "L", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(widths[0]+widths[1]+widths[2]+widths[3]+widths[4], 8, "TOTALS", "1", 0, "L", false, 0, "")
	pdf.CellFormat(widths[5], 8, fmt.Sprintf("%.2f", totalAmount), "1", 0, "R", false, 0, "")
	pdf.CellFormat(widths[6], 8, "", "1", 0, "L", false, 0, "")
	pdf.CellFormat(widths[7], 8, "", "1", 0, "L", false, 0, "")
	pdf.CellFormat(widths[8], 8, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(widths[9], 8, "", "1", 0, "L", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}
