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

// GetMemberPayrollDeductions retrieves a list of member payroll deductions with optional filtering and pagination.
func (s *MemberPayrollDeductionService) GetMemberPayrollDeductions(page, limit int, memberID, payrollID, deductionMonth, fiscalYear, settled, confirmed string) ([]dtos.MemberPayrollDeductionResponse, int64, error) {
	var deductions []dtos.MemberPayrollDeductionResponse
	var total int64

	query := db.DB.Table("member_payroll_deductions mpd").
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

	if memberID != "" {
		query = query.Where("mpd.member_id = ?", memberID)
	}
	if payrollID != "" {
		query = query.Where("mpd.payroll_id = ?", payrollID)
	}
	if deductionMonth != "" {
		query = query.Where("mpd.deduction_month = ?", deductionMonth)
	}
	if fiscalYear != "" {
		query = query.Where("mpd.fiscal_year = ?", fiscalYear)
	}
	if settled != "" {
		query = query.Where("mpd.settled = ?", settled)
	}
	if confirmed != "" {
		query = query.Where("mpd.confirmed = ?", confirmed)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Order("mpd.id DESC").Scan(&deductions).Error

	return deductions, total, err
}

// GetMemberPayrollDeduction retrieves a single member payroll deduction by its ID.
func (s *MemberPayrollDeductionService) GetMemberPayrollDeduction(id string) (*dtos.MemberPayrollDeductionResponse, error) {
	var deduction dtos.MemberPayrollDeductionResponse
	query := db.DB.Table("member_payroll_deductions mpd").
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
		Where("mpd.id = ? AND mpd.deleted_at IS NULL", id).
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
func (s *MemberPayrollDeductionService) ExportDeductions(userID uint64, memberID, payrollID, deductionMonth, fiscalYear, settled, confirmed, reportType string) error {
	go s.processDeductionExportInBackground(userID, memberID, payrollID, deductionMonth, fiscalYear, settled, confirmed, reportType)
	return nil
}

func (s *MemberPayrollDeductionService) processDeductionExportInBackground(userID uint64, memberID, payrollID, deductionMonth, fiscalYear, settled, confirmed, reportType string) {
	var results []dtos.MemberPayrollDeductionResponse

	baseQuery := db.DB.Table("member_payroll_deductions mpd").
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

	if memberID != "" {
		baseQuery = baseQuery.Where("mpd.member_id = ?", memberID)
	}
	if payrollID != "" {
		baseQuery = baseQuery.Where("mpd.payroll_id = ?", payrollID)
	}
	if deductionMonth != "" {
		baseQuery = baseQuery.Where("mpd.deduction_month = ?", deductionMonth)
	}
	if fiscalYear != "" {
		baseQuery = baseQuery.Where("mpd.fiscal_year = ?", fiscalYear)
	}
	if settled != "" {
		baseQuery = baseQuery.Where("mpd.settled = ?", settled)
	}
	if confirmed != "" {
		baseQuery = baseQuery.Where("mpd.confirmed = ?", confirmed)
	}

	if err := baseQuery.Order("mpd.id DESC").Scan(&results).Error; err != nil {
		log.Printf("[MemberPayrollDeductionService.processDeductionExportInBackground] Query error: %v", err)
		return
	}

	var fileData []byte
	var err error
	ext := "csv"
	if strings.ToLower(reportType) == "pdf" {
		ext = "pdf"
		fileData, err = s.generateDeductionPDF(results)
	} else {
		fileData, err = s.generateDeductionCSV(results)
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

func (s *MemberPayrollDeductionService) generateDeductionCSV(results []dtos.MemberPayrollDeductionResponse) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	headers := []string{"Member No", "Member Name", "Month", "Year", "Type", "Amount", "Payroll Period", "Reference", "Settled"}
	writer.Write(headers)

	var totalAmount float64

	for _, d := range results {
		amt, _ := strconv.ParseFloat(d.Amount, 64)
		totalAmount += amt

		writer.Write([]string{
			d.MemberNo,
			d.MemberName,
			d.DeductionMonth,
			fmt.Sprintf("%d", d.FiscalYear),
			d.DeductionTypeName,
			d.Amount,
			d.FiscalPeriod,
			d.Reference,
			d.Settled,
		})
	}

	writer.Write([]string{
		"TOTALS", "", "", "", "",
		fmt.Sprintf("%.2f", totalAmount),
		"", "", "",
	})

	writer.Flush()
	return buf.Bytes(), writer.Error()
}

func (s *MemberPayrollDeductionService) generateDeductionPDF(results []dtos.MemberPayrollDeductionResponse) ([]byte, error) {
	var org struct {
		RegisteredName string `gorm:"column:registered_name"`
		Address        string `gorm:"column:address"`
		Phone          string `gorm:"column:phone"`
		Email          string `gorm:"column:email"`
	}
	db.DB.Table("organization_details").First(&org)

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, org.RegisteredName, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, org.Address, "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Phone: %s | Email: %s", org.Phone, org.Email), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "MEMBER PAYROLL DEDUCTIONS REGISTER", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 8)
	headers := []string{"M-No", "Member Name", "Month", "Year", "Type", "Amount", "Payroll Period", "Reference", "Settled"}
	widths := []float64{25, 55, 20, 15, 40, 25, 25, 50, 20}

	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	var totalAmount float64

	pdf.SetFont("Arial", "", 8)
	for _, d := range results {
		amt, _ := strconv.ParseFloat(d.Amount, 64)
		totalAmount += amt

		pdf.CellFormat(widths[0], 8, d.MemberNo, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 8, d.MemberName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 8, d.DeductionMonth, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 8, fmt.Sprintf("%d", d.FiscalYear), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[4], 8, d.DeductionTypeName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[5], 8, d.Amount, "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[6], 8, d.FiscalPeriod, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[7], 8, d.Reference, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[8], 8, d.Settled, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(widths[0], 8, "TOTALS", "1", 0, "L", false, 0, "")
	pdf.CellFormat(widths[1], 8, "", "1", 0, "L", false, 0, "")
	pdf.CellFormat(widths[2], 8, "", "1", 0, "L", false, 0, "")
	pdf.CellFormat(widths[3], 8, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(widths[4], 8, "", "1", 0, "L", false, 0, "")
	pdf.CellFormat(widths[5], 8, fmt.Sprintf("%.2f", totalAmount), "1", 0, "R", false, 0, "")
	pdf.CellFormat(widths[6], 8, "", "1", 0, "L", false, 0, "")
	pdf.CellFormat(widths[7], 8, "", "1", 0, "L", false, 0, "")
	pdf.CellFormat(widths[8], 8, "", "1", 0, "C", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}
