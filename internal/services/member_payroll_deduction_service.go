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
			mpd.confirmed, mpd.payroll_id, mp.fiscal_period, mpd.reference, mpd.settlement_type,
			mpd.created_at, mpd.updated_at
		`).
		Joins("LEFT JOIN member_registrations mr ON mpd.member_id = mr.id").
		Joins("LEFT JOIN deduction_types dt ON mpd.deduction_type_id = dt.id").
		Joins("LEFT JOIN member_payrolls mp ON mpd.payroll_id = mp.id").
		Where("mpd.deleted_at IS NULL")
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

	return deductions, total, err
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
	return &deduction, nil
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
