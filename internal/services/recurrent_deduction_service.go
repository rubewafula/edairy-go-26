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

// RecurrentDeductionFilters holds shared list/export filter values.
type RecurrentDeductionFilters struct {
	CustomerID      string
	CustomerType    string
	Settled         string
	MemberNo        string
	DeductionTypeID string
	DateFrom        string
	DateTo          string
}

func (s *RecurrentDeductionService) baseRecurrentDeductionQuery() *gorm.DB {
	return db.DB.Table("recurrent_deductions rd").
		Select(`
			rd.*,
			mr.member_no,
			CONCAT(mr.first_name, ' ', mr.last_name) as names,
			dt.description as deduction_type_name
		`).
		Joins("LEFT JOIN member_registrations mr ON rd.customer_id = mr.id AND rd.customer_type = 'member'").
		Joins("LEFT JOIN deduction_types dt ON rd.deduction_type_id = dt.id").
		Where("rd.deleted_at IS NULL")
}

func (s *RecurrentDeductionService) applyRecurrentDeductionFilters(query *gorm.DB, filters RecurrentDeductionFilters) *gorm.DB {
	if filters.CustomerID != "" {
		query = query.Where("rd.customer_id = ?", filters.CustomerID)
	}
	if filters.CustomerType != "" {
		query = query.Where("rd.customer_type = ?", filters.CustomerType)
	}
	if filters.Settled != "" {
		query = query.Where("rd.settled = ?", filters.Settled)
	}
	if filters.MemberNo != "" {
		query = query.Where("mr.member_no LIKE ?", "%"+filters.MemberNo+"%")
	}
	if filters.DeductionTypeID != "" {
		query = query.Where("rd.deduction_type_id = ?", filters.DeductionTypeID)
	}
	if filters.DateFrom != "" {
		query = query.Where("DATE(rd.transaction_date) >= ?", filters.DateFrom)
	}
	if filters.DateTo != "" {
		query = query.Where("DATE(rd.transaction_date) <= ?", filters.DateTo)
	}
	return query
}

// GetRecurrentDeductions retrieves a list of recurrent deductions with optional filtering and pagination.
func (s *RecurrentDeductionService) GetRecurrentDeductions(page, limit int, filters RecurrentDeductionFilters) ([]dtos.RecurrentDeductionResponse, int64, error) {
	var deductions []dtos.RecurrentDeductionResponse
	var total int64

	query := s.applyRecurrentDeductionFilters(s.baseRecurrentDeductionQuery(), filters)

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Order("rd.id DESC").Scan(&deductions).Error

	return deductions, total, err
}

// GetRecurrentDeduction retrieves a single recurrent deduction by its ID.
func (s *RecurrentDeductionService) GetRecurrentDeduction(id string) (*dtos.RecurrentDeductionResponse, error) {
	var deduction dtos.RecurrentDeductionResponse
	query := s.baseRecurrentDeductionQuery().
		Where("rd.id = ?", id).
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
func (s *RecurrentDeductionService) ExportRecurrentDeductions(userID uint64, filters RecurrentDeductionFilters, reportType string) error {
	go s.processRecurrentDeductionExportInBackground(userID, filters, reportType)
	return nil
}

func (s *RecurrentDeductionService) processRecurrentDeductionExportInBackground(userID uint64, filters RecurrentDeductionFilters, reportType string) {
	var results []dtos.RecurrentDeductionResponse

	query := s.applyRecurrentDeductionFilters(s.baseRecurrentDeductionQuery(), filters)

	if err := query.Order("rd.id DESC").Scan(&results).Error; err != nil {
		log.Printf("[RecurrentDeductionService.processRecurrentDeductionExportInBackground] Query error: %v", err)
		return
	}

	title := s.buildRecurrentDeductionExportTitle(filters, results)

	var fileData []byte
	var err error
	ext := "csv"
	if strings.ToLower(reportType) == "pdf" {
		ext = "pdf"
		fileData, err = s.generateRecurrentDeductionPDF(results, title)
	} else {
		fileData, err = s.generateRecurrentDeductionCSV(results, title)
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

type recurrentDeductionExportTitle struct {
	DairyName    string
	PlatformName string
	ReportTitle  string
}

func (s *RecurrentDeductionService) formatExportDate(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if t := utils.ParseDate(value); !t.IsZero() {
		return t.Format("02 Jan 2006")
	}
	return value
}

func (s *RecurrentDeductionService) resolveExportMemberLabel(filters RecurrentDeductionFilters, results []dtos.RecurrentDeductionResponse) string {
	if filters.CustomerID == "" && strings.TrimSpace(filters.MemberNo) == "" {
		return ""
	}

	for _, row := range results {
		memberNo := strings.TrimSpace(row.MemberNo)
		names := strings.TrimSpace(row.Names)
		if memberNo == "" && names == "" {
			continue
		}
		if memberNo != "" && names != "" {
			return strings.TrimSpace(memberNo + " " + names)
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
	if filters.CustomerID != "" {
		query = query.Where("id = ?", filters.CustomerID)
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

func (s *RecurrentDeductionService) resolveExportDeductionTypeName(filters RecurrentDeductionFilters, results []dtos.RecurrentDeductionResponse) string {
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

func (s *RecurrentDeductionService) buildRecurrentDeductionExportTitle(filters RecurrentDeductionFilters, results []dtos.RecurrentDeductionResponse) recurrentDeductionExportTitle {
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

	platformName := "Edairy Africa Ltd"

	reportTitle := "Member Deductions"
	fromLabel := s.formatExportDate(filters.DateFrom)
	toLabel := s.formatExportDate(filters.DateTo)
	switch {
	case fromLabel != "" && toLabel != "":
		reportTitle = fmt.Sprintf("Member Deductions From %s To %s", fromLabel, toLabel)
	case fromLabel != "":
		reportTitle = fmt.Sprintf("Member Deductions From %s", fromLabel)
	case toLabel != "":
		reportTitle = fmt.Sprintf("Member Deductions To %s", toLabel)
	}

	if memberLabel := s.resolveExportMemberLabel(filters, results); memberLabel != "" {
		reportTitle = fmt.Sprintf("%s - Member %s", reportTitle, memberLabel)
	}
	if typeName := s.resolveExportDeductionTypeName(filters, results); typeName != "" {
		reportTitle = fmt.Sprintf("%s - %s", reportTitle, typeName)
	}

	return recurrentDeductionExportTitle{
		DairyName:    dairyName,
		PlatformName: platformName,
		ReportTitle:  reportTitle,
	}
}

func (s *RecurrentDeductionService) generateRecurrentDeductionCSV(results []dtos.RecurrentDeductionResponse, title recurrentDeductionExportTitle) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	_ = writer.Write([]string{title.DairyName})
	_ = writer.Write([]string{title.PlatformName})
	_ = writer.Write([]string{title.ReportTitle})
	_ = writer.Write([]string{})

	headers := []string{"Customer ID", "Customer Type", "Member No", "Names", "Deduction Type", "Reference", "Total Amount", "Paid Amount", "Recurrent Amount", "Principal Amount", "Settled", "Transaction Date"}
	_ = writer.Write(headers)

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

		_ = writer.Write([]string{
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

	_ = writer.Write([]string{
		"TOTALS", "", "", "", "", "",
		fmt.Sprintf("%.2f", totalTotalAmount),
		fmt.Sprintf("%.2f", totalPaidAmount),
		fmt.Sprintf("%.2f", totalRecurrentAmount),
		fmt.Sprintf("%.2f", totalPrincipalAmount),
		"", "",
	})
	_ = writer.Write([]string{})
	_ = writer.Write([]string{"powered by edairy africa"})

	writer.Flush()
	return buf.Bytes(), writer.Error()
}

func (s *RecurrentDeductionService) generateRecurrentDeductionPDF(results []dtos.RecurrentDeductionResponse, title recurrentDeductionExportTitle) ([]byte, error) {
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
	headers := []string{"Cust ID", "Cust Type", "M-No", "Names", "Deduction Type", "Reference", "Total Amt", "Paid Amt", "Recurrent Amt", "Principal Amt", "Settled", "Trans Date"}
	widths := []float64{15, 20, 20, 40, 30, 30, 20, 20, 20, 20, 15, 20}

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
