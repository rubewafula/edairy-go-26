package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type EmployeeLeaveApplicationService struct {
	notificationService *UINotificationService
}

func NewEmployeeLeaveApplicationService() *EmployeeLeaveApplicationService {
	return &EmployeeLeaveApplicationService{
		notificationService: NewUINotificationService(),
	}
}

func (s *EmployeeLeaveApplicationService) CreateApplication(req dtos.CreateEmployeeLeaveApplicationRequest, userID uint64) (*models.EmployeeLeaveApplication, error) {
	application := &models.EmployeeLeaveApplication{
		BaseModel:     models.BaseModel{CreatedBy: userID},
		ApplicationNo: fmt.Sprintf("LV-%d-%d", req.EmployeeID, time.Now().Unix()),
		EmployeeID:    req.EmployeeID,
		LeaveTypeID:   req.LeaveTypeID,
		DaysApplied:   req.DaysApplied,
		StartDate:     utils.ParseDate(req.StartDate),
		EndDate:       utils.ParseDate(req.EndDate),
		ReturnDate:    utils.ParseDate(req.ReturnDate),
		Status:        "PENDING",
	}

	if err := db.DB.Create(application).Error; err != nil {
		return nil, err
	}
	return application, nil
}

func (s *EmployeeLeaveApplicationService) GetApplications(employeeID string, page, limit int) ([]dtos.EmployeeLeaveApplicationResponse, int64, error) {
	var results []dtos.EmployeeLeaveApplicationResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeeLeaveApplication{})
	if employeeID != "" {
		queryBuilder = queryBuilder.Where("employee_id = ?", employeeID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			ela.*, 
			CONCAT(e.first_name, ' ', e.surname) as employee_name,
			elt.description as leave_type,
			CONCAT(app.first_name, ' ', app.surname) as approver_name
		FROM employee_leave_applications ela
		LEFT JOIN employees e ON ela.employee_id = e.id
		LEFT JOIN employee_leave_types elt ON ela.leave_type_id = elt.id
		LEFT JOIN employees app ON ela.approver_id = app.id
		WHERE ela.deleted_at IS NULL AND (? = '' OR ela.employee_id = ?)
		ORDER BY ela.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, employeeID, employeeID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeeLeaveApplicationService) GetApplication(id string) (*dtos.EmployeeLeaveApplicationResponse, error) {
	var result dtos.EmployeeLeaveApplicationResponse
	query := `
		SELECT 
			ela.*, 
			CONCAT(e.first_name, ' ', e.surname) as employee_name,
			elt.description as leave_type,
			CONCAT(app.first_name, ' ', app.surname) as approver_name
		FROM employee_leave_applications ela
		LEFT JOIN employees e ON ela.employee_id = e.id
		LEFT JOIN employee_leave_types elt ON ela.leave_type_id = elt.id
		LEFT JOIN employees app ON ela.approver_id = app.id
		WHERE ela.id = ? AND ela.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *EmployeeLeaveApplicationService) UpdateApplication(id string, req dtos.UpdateEmployeeLeaveApplicationRequest, userID uint64) error {
	var application models.EmployeeLeaveApplication
	if err := db.DB.First(&application, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"approver_id":   req.ApproverID,
		"days_approved": req.DaysApproved,
		"status":        req.Status,
		"approved":      req.Approved,
		"updated_by":    userID,
	}

	return db.DB.Model(&application).Updates(updates).Error
}

func (s *EmployeeLeaveApplicationService) DeleteApplication(id string, userID uint64) error {
	var application models.EmployeeLeaveApplication
	if err := db.DB.First(&application, id).Error; err != nil {
		return err
	}
	// Audit before soft delete
	return db.DB.Model(&application).Update("updated_by", userID).Delete(&application).Error
}

// ExportApplications initiates a background process to export leave applications.
func (s *EmployeeLeaveApplicationService) ExportApplications(userID uint64, status, format string) error {
	go s.processExportInBackground(userID, status, format)
	return nil
}

type leaveApplicationExportResult struct {
	EmployeeNo       string    `gorm:"column:employee_no"`
	EmployeeNames    string    `gorm:"column:employee_names"`
	Position         string    `gorm:"column:position"`
	LeaveCode        string    `gorm:"column:leave_code"`
	LeaveDescription string    `gorm:"column:leave_description"`
	StartDate        time.Time `gorm:"column:start_date"`
	EndDate          time.Time `gorm:"column:end_date"`
	DaysApproved     float64   `gorm:"column:days_approved"`
	BalanceBF        string    `gorm:"column:balance_bf"`
}

func (s *EmployeeLeaveApplicationService) processExportInBackground(userID uint64, status, format string) {
	var results []leaveApplicationExportResult

	query := `
		SELECT 
			e.employee_no, 
			CONCAT(e.first_name, ' ', COALESCE(e.middle_name, ''), ' ', e.surname) as employee_names,
			jp.name as position,
			elt.code as leave_code,
			elt.description as leave_description,
			ela.start_date,
			ela.end_date,
			ela.days_approved,
			eld.balance_bf
		FROM employee_leave_applications ela
		JOIN employees e ON ela.employee_id = e.id
		LEFT JOIN job_positions jp ON e.job_position_id = jp.id
		JOIN employee_leave_types elt ON ela.leave_type_id = elt.id
		LEFT JOIN employee_leave_details eld ON ela.employee_id = eld.employee_id
		WHERE ela.deleted_at IS NULL`

	if status != "" {
		query += fmt.Sprintf(" AND ela.status = '%s'", status)
	}

	if err := db.DB.Raw(query).Scan(&results).Error; err != nil {
		log.Printf("[EmployeeLeaveApplicationService] Export query error: %v", err)
		return
	}

	var fileData []byte
	var err error
	ext := "csv"

	if strings.ToLower(format) == "pdf" {
		ext = "pdf"
		fileData, err = s.generatePDF(results)
	} else {
		buf := new(bytes.Buffer)
		writer := csv.NewWriter(buf)
		writer.Write([]string{"Employee No", "Employee Names", "Position", "leave_code", "leave_description", "start_date", "end_date", "days_approved", "balance_bf"})

		for _, r := range results {
			writer.Write([]string{
				r.EmployeeNo,
				r.EmployeeNames,
				r.Position,
				r.LeaveCode,
				r.LeaveDescription,
				r.StartDate.Format("2006-01-02"),
				r.EndDate.Format("2006-01-02"),
				fmt.Sprintf("%.1f", r.DaysApproved),
				r.BalanceBF,
			})
		}
		writer.Flush()
		fileData = buf.Bytes()
		err = writer.Error()
	}

	if err != nil {
		log.Printf("[EmployeeLeaveApplicationService] Error generating export: %v", err)
		return
	}

	exportDir := "./storage/exports"
	os.MkdirAll(exportDir, 0755)
	filename := fmt.Sprintf("leave_applications_%d.%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(exportDir, filename)

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		log.Printf("[EmployeeLeaveApplicationService] Error saving export file: %v", err)
		return
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            fmt.Sprintf("Leave Applications Export (%s) Ready", strings.ToUpper(ext)),
		Message:          fmt.Sprintf("The leave applications %s export is ready for download.", ext),
		NotificationType: "SUCCESS",
		DownloadLink:     fmt.Sprintf("/api/employee-leave-applications/export/download/%s", filename),
	})
}

func (s *EmployeeLeaveApplicationService) generatePDF(results []leaveApplicationExportResult) ([]byte, error) {
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
	pdf.CellFormat(0, 10, "EMPLOYEE LEAVE APPLICATIONS REGISTER", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 8)
	headers := []string{"Emp-No", "Names", "Position", "Code", "Type", "Start", "End", "Days", "Bal BF"}
	widths := []float64{20, 55, 45, 20, 45, 25, 25, 20, 20}

	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 8)
	for _, r := range results {
		pdf.CellFormat(widths[0], 8, r.EmployeeNo, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 8, r.EmployeeNames, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 8, r.Position, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 8, r.LeaveCode, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[4], 8, r.LeaveDescription, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[5], 8, r.StartDate.Format("2006-01-02"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[6], 8, r.EndDate.Format("2006-01-02"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[7], 8, fmt.Sprintf("%.1f", r.DaysApproved), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[8], 8, r.BalanceBF, "1", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}
