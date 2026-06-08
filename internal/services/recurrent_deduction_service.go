package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type RecurrentDeductionService struct {
	notificationService *UINotificationService
}

func NewRecurrentDeductionService() *RecurrentDeductionService {
	return &RecurrentDeductionService{
		notificationService: NewUINotificationService(),
	}
}

// GetRecurrentDeductions retrieves a list of recurrent deductions with optional filtering and pagination.
func (s *RecurrentDeductionService) GetRecurrentDeductions(page, limit int, customerID string, customerType string, settled string) ([]dtos.RecurrentDeductionResponse, int64, error) {
	var deductions []dtos.RecurrentDeductionResponse
	var total int64

	query := db.DB.Table("recurrent_deductions rd").
		Select(`
			rd.*,
			mr.member_no,
			CONCAT(mr.first_name, ' ', mr.last_name) as names,
			dt.description as deduction_type_name
		`).
		Joins("LEFT JOIN member_registrations mr ON rd.customer_id = mr.id AND rd.customer_type = 'member'").
		Joins("LEFT JOIN deduction_types dt ON rd.deduction_type_id = dt.id").
		Where("rd.deleted_at IS NULL")

	if customerID != "" {
		query = query.Where("rd.customer_id = ?", customerID)
	}
	if customerType != "" {
		query = query.Where("rd.customer_type = ?", customerType)
	}
	if settled != "" {
		query = query.Where("rd.settled = ?", settled)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Order("rd.id DESC").Scan(&deductions).Error

	return deductions, total, err
}

// GetRecurrentDeduction retrieves a single recurrent deduction by its ID.
func (s *RecurrentDeductionService) GetRecurrentDeduction(id string) (*dtos.RecurrentDeductionResponse, error) {
	var deduction dtos.RecurrentDeductionResponse
	query := db.DB.Table("recurrent_deductions rd").
		Select(`
			rd.*,
			mr.member_no,
			CONCAT(mr.first_name, ' ', mr.last_name) as names,
			dt.description as deduction_type_name
		`).
		Joins("LEFT JOIN member_registrations mr ON rd.customer_id = mr.id AND rd.customer_type = 'member'").
		Joins("LEFT JOIN deduction_types dt ON rd.deduction_type_id = dt.id").
		Where("rd.id = ? AND rd.deleted_at IS NULL", id).
		Limit(1)

	if err := query.Scan(&deduction).Error; err != nil {
		return nil, err
	}

	if deduction.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &deduction, nil
}

// ExportRecurrentDeductions initiates a background process to export recurrent deductions.
func (s *RecurrentDeductionService) ExportRecurrentDeductions(userID uint64, customerID, customerType, settled, reportType string) error {
	go s.processRecurrentDeductionExportInBackground(userID, customerID, customerType, settled, reportType)
	return nil
}

func (s *RecurrentDeductionService) processRecurrentDeductionExportInBackground(userID uint64, customerID, customerType, settled, reportType string) {
	var results []dtos.RecurrentDeductionResponse

	query := db.DB.Table("recurrent_deductions rd").
		Select(`
			rd.*,
			mr.member_no,
			CONCAT(mr.first_name, ' ', mr.last_name) as names,
			dt.description as deduction_type_name
		`).
		Joins("LEFT JOIN member_registrations mr ON rd.customer_id = mr.id AND rd.customer_type = 'member'").
		Joins("LEFT JOIN deduction_types dt ON rd.deduction_type_id = dt.id").
		Where("rd.deleted_at IS NULL")

	if customerID != "" {
		query = query.Where("rd.customer_id = ?", customerID)
	}
	if customerType != "" {
		query = query.Where("rd.customer_type = ?", customerType)
	}
	if settled != "" {
		query = query.Where("rd.settled = ?", settled)
	}

	if err := query.Order("rd.id DESC").Scan(&results).Error; err != nil {
		log.Printf("[RecurrentDeductionService.processRecurrentDeductionExportInBackground] Query error: %v", err)
		return
	}

	var fileData []byte
	var err error
	ext := "csv"
	if strings.ToLower(reportType) == "pdf" {
		ext = "pdf"
		fileData, err = s.generateRecurrentDeductionPDF(results, customerID, customerType, settled)
	} else {
		fileData, err = s.generateRecurrentDeductionCSV(results)
	}

	if err != nil {
		log.Printf("[RecurrentDeductionService.processRecurrentDeductionExportInBackground] Generation error: %v", err)
		return
	}

	exportDir := "./storage/exports"
	os.MkdirAll(exportDir, 0755)
	filename := fmt.Sprintf("recurrent_deductions_export_%d.%s", utils.Now().UnixNano(), ext)
	filePath := filepath.Join(exportDir, filename)

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		log.Printf("[RecurrentDeductionService.processRecurrentDeductionExportInBackground] File error: %v", err)
		return
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Recurrent Deductions Export Complete",
		Message:          fmt.Sprintf("Your recurrent deductions export (%s) is ready for download.", strings.ToUpper(ext)),
		NotificationType: "SUCCESS",
		DownloadLink:     fmt.Sprintf("/api/member-deductions/export/download/%s", filename),
	})
}

func (s *RecurrentDeductionService) generateRecurrentDeductionCSV(results []dtos.RecurrentDeductionResponse) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	headers := []string{"Customer ID", "Customer Type", "Member No", "Names", "Deduction Type", "Reference", "Total Amount", "Paid Amount", "Recurrent Amount", "Principal Amount", "Settled", "Transaction Date"}
	writer.Write(headers)

	var totalTotalAmount, totalPaidAmount, totalRecurrentAmount, totalPrincipalAmount float64

	for _, rd := range results {
		totalTotalAmount += rd.TotalAmount
		totalPaidAmount += rd.PaidAmount
		totalRecurrentAmount += rd.RecurrentAmount
		totalPrincipalAmount += rd.PrincipalAmount

		transDate := ""
		if rd.TransactionDate != nil {
			transDate = rd.TransactionDate.Format("2006-01-02")
		}

		writer.Write([]string{
			fmt.Sprintf("%d", rd.CustomerID),
			rd.CustomerType,
			rd.MemberNo,
			rd.Names,
			rd.DeductionTypeName,
			rd.Reference,
			fmt.Sprintf("%.2f", rd.TotalAmount),
			fmt.Sprintf("%.2f", rd.PaidAmount),
			fmt.Sprintf("%.2f", rd.RecurrentAmount),
			fmt.Sprintf("%.2f", rd.PrincipalAmount),
			fmt.Sprintf("%d", rd.Settled),
			transDate,
		})
	}

	writer.Write([]string{
		"TOTALS", "", "", "", "", "",
		fmt.Sprintf("%.2f", totalTotalAmount),
		fmt.Sprintf("%.2f", totalPaidAmount),
		fmt.Sprintf("%.2f", totalRecurrentAmount),
		fmt.Sprintf("%.2f", totalPrincipalAmount),
		"", "",
	})

	writer.Flush()
	return buf.Bytes(), writer.Error()
}

func (s *RecurrentDeductionService) generateRecurrentDeductionPDF(results []dtos.RecurrentDeductionResponse, customerID, customerType, settled string) ([]byte, error) {
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
	pdf.CellFormat(0, 10, "RECURRENT DEDUCTIONS REGISTER", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Filter Summary
	filterSummary := "Filters: "
	if customerID != "" {
		filterSummary += fmt.Sprintf("Customer ID: %s, ", customerID)
	}
	if customerType != "" {
		filterSummary += fmt.Sprintf("Customer Type: %s, ", customerType)
	}
	if settled != "" {
		filterSummary += fmt.Sprintf("Settled: %s, ", settled)
	}
	if len(filterSummary) > len("Filters: ") {
		filterSummary = strings.TrimSuffix(filterSummary, ", ")
		pdf.SetFont("Arial", "I", 9)
		pdf.CellFormat(0, 8, filterSummary, "", 1, "L", false, 0, "")
		pdf.Ln(2)
	}

	pdf.SetFont("Arial", "B", 8)
	headers := []string{"Cust ID", "Cust Type", "M-No", "Names", "Deduction Type", "Reference", "Total Amt", "Paid Amt", "Recurrent Amt", "Principal Amt", "Settled", "Trans Date"}
	widths := []float64{15, 20, 20, 40, 30, 30, 20, 20, 20, 20, 15, 20} // Adjusted widths for L (landscape) A4

	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	var totalTotalAmount, totalPaidAmount, totalRecurrentAmount, totalPrincipalAmount float64

	pdf.SetFont("Arial", "", 8)
	for _, rd := range results {
		totalTotalAmount += rd.TotalAmount
		totalPaidAmount += rd.PaidAmount
		totalRecurrentAmount += rd.RecurrentAmount
		totalPrincipalAmount += rd.PrincipalAmount

		transDate := ""
		if rd.TransactionDate != nil {
			transDate = rd.TransactionDate.Format("2006-01-02")
		}

		pdf.CellFormat(widths[0], 8, fmt.Sprintf("%d", rd.CustomerID), "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 8, rd.CustomerType, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 8, rd.MemberNo, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 8, rd.Names, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[4], 8, rd.DeductionTypeName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[5], 8, rd.Reference, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[6], 8, fmt.Sprintf("%.2f", rd.TotalAmount), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[7], 8, fmt.Sprintf("%.2f", rd.PaidAmount), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[8], 8, fmt.Sprintf("%.2f", rd.RecurrentAmount), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[9], 8, fmt.Sprintf("%.2f", rd.PrincipalAmount), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[10], 8, fmt.Sprintf("%d", rd.Settled), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[11], 8, transDate, "1", 0, "L", false, 0, "")
		pdf.Ln(-1)
	}

	// Totals row
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(widths[0]+widths[1]+widths[2]+widths[3]+widths[4]+widths[5], 8, "TOTALS", "1", 0, "L", false, 0, "")
	pdf.CellFormat(widths[6], 8, fmt.Sprintf("%.2f", totalTotalAmount), "1", 0, "R", false, 0, "")
	pdf.CellFormat(widths[7], 8, fmt.Sprintf("%.2f", totalPaidAmount), "1", 0, "R", false, 0, "")
	pdf.CellFormat(widths[8], 8, fmt.Sprintf("%.2f", totalRecurrentAmount), "1", 0, "R", false, 0, "")
	pdf.CellFormat(widths[9], 8, fmt.Sprintf("%.2f", totalPrincipalAmount), "1", 0, "R", false, 0, "")
	pdf.CellFormat(widths[10], 8, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(widths[11], 8, "", "1", 0, "L", false, 0, "")
	pdf.Ln(-1)

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}
