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

type EmployeeService struct {
	notificationService *UINotificationService
}

func NewEmployeeService() *EmployeeService {
	return &EmployeeService{
		notificationService: NewUINotificationService(),
	}
}

func (s *EmployeeService) CreateEmployee(req dtos.CreateEmployeeRequest, userID uint64) (*models.Employee, error) {
	employee := &models.Employee{
		BaseModel:         models.BaseModel{CreatedBy: userID},
		UserID:            req.UserID,
		Surname:           req.Surname,
		FirstName:         req.FirstName,
		MiddleName:        req.MiddleName,
		EmployeeNo:        req.EmployeeNo,
		IDNo:              req.IDNo,
		KraPin:            req.KraPin,
		NssfNo:            req.NssfNo,
		NhifNo:            req.NhifNo,
		Gender:            req.Gender,
		DateOfBirth:       utils.ParseDate(req.DateOfBirth),
		Phone:             req.Phone,
		Email:             req.Email,
		JobPositionID:     req.JobPositionID,
		Status:            req.Status,
		Title:             req.Title,
		Town:              req.Town,
		SiteID:            req.SiteID,
		MaritalStatus:     req.MaritalStatus,
		Religion:          req.Religion,
		Disabled:          req.Disabled,
		StoreID:           req.StoreID,
		PostalAddress:     req.PostalAddress,
		PostalCode:        req.PostalCode,
		BirthCity:         req.BirthCity,
		NextOfKinFullName: req.NextOfKinFullName,
		NextOfKinPhone:    req.NextOfKinPhone,
	}

	if err := db.DB.Create(employee).Error; err != nil {
		return nil, err
	}
	return employee, nil
}

func (s *EmployeeService) GetEmployees(page, limit int) ([]dtos.EmployeeResponse, int64, error) {
	var employees []dtos.EmployeeResponse
	var total int64
	db.DB.Model(&models.Employee{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT e.*, jp.name as job_position_name, d.department_name as department_name
		FROM employees e
		LEFT JOIN job_positions jp ON e.job_position_id = jp.id
		LEFT JOIN departments d ON jp.department_id = d.id
		WHERE e.deleted_at IS NULL
		ORDER BY e.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&employees).Error
	return employees, total, err
}

func (s *EmployeeService) GetEmployee(id string) (*dtos.EmployeeResponse, error) {
	var employee dtos.EmployeeResponse
	query := `
		SELECT e.*, jp.name as job_position_name, d.department_name as department_name
		FROM employees e
		LEFT JOIN job_positions jp ON e.job_position_id = jp.id
		LEFT JOIN departments d ON jp.department_id = d.id
		WHERE e.id = ? AND e.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&employee).Error
	if err != nil {
		return nil, err
	}
	if employee.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &employee, nil
}

func (s *EmployeeService) UpdateEmployee(id string, req dtos.UpdateEmployeeRequest, userID uint64) error {
	var employee models.Employee
	if err := db.DB.First(&employee, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"surname":       req.Surname,
		"first_name":    req.FirstName,
		"middle_name":   req.MiddleName,
		"phone_number":  req.Phone,
		"email_address": req.Email,
		"updated_by":    userID,
	}

	return db.DB.Model(&employee).Updates(updates).Error
}

func (s *EmployeeService) DeleteEmployee(id string, userID uint64) error {
	var employee models.Employee
	if err := db.DB.First(&employee, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&employee).Update("updated_by", userID).Delete(&employee).Error
}

func (s *EmployeeService) GetSalaries(employeeID string) ([]models.EmployeeSalary, error) {
	var salaries []models.EmployeeSalary
	err := db.DB.Where("employee_id = ?", employeeID).Find(&salaries).Error
	return salaries, err
}

func (s *EmployeeService) GetSalary(id string) (*models.EmployeeSalary, error) {
	var salary models.EmployeeSalary
	if err := db.DB.First(&salary, id).Error; err != nil {
		return nil, err
	}
	return &salary, nil
}

func (s *EmployeeService) UpdateSalary(id string, req dtos.UpdateEmployeeSalaryRequest, userID uint64) error {
	return db.DB.Model(&models.EmployeeSalary{}).Where("id = ?", id).
		Updates(map[string]interface{}{"basic_salary": req.BasicSalary, "status": req.Status, "updated_by": userID}).Error
}

func (s *EmployeeService) CreateSalary(req dtos.CreateEmployeeSalaryRequest, userID uint64) (*models.EmployeeSalary, error) {
	salary := &models.EmployeeSalary{
		BaseModel:   models.BaseModel{CreatedBy: userID},
		EmployeeID:  req.EmployeeID,
		BasicSalary: req.BasicSalary,
		Status:      req.Status,
	}
	if salary.Status == "" {
		salary.Status = "ACTIVE"
	}
	err := db.DB.Create(salary).Error
	return salary, err
}

// ExportEmployees initiates a background process to export employee data to CSV or PDF.
func (s *EmployeeService) ExportEmployees(userID uint64, status, format string) error {
	go s.processEmployeeExportInBackground(userID, status, format)
	return nil
}

type employeeExportQueryResult struct {
	EmployeeNo    string     `gorm:"column:employee_no"`
	IdNo          string     `gorm:"column:id_no"`
	FullName      string     `gorm:"column:full_name"`
	DateOfBirth   *time.Time `gorm:"column:date_of_birth"`
	Gender        string     `gorm:"column:gender"`
	MaritalStatus string     `gorm:"column:marital_status"`
	PositionName  string     `gorm:"column:position_name"`
}

func (s *EmployeeService) processEmployeeExportInBackground(userID uint64, status, format string) {
	var results []employeeExportQueryResult

	// Query to fetch employee details with job position join for requested columns
	query := `
		SELECT 
			e.employee_no, 
			e.id_no, 
			CONCAT(e.first_name, ' ', COALESCE(e.middle_name, ''), ' ', e.surname) as full_name,
			e.date_of_birth, 
			e.gender, 
			e.marital_status, 
			jp.name as position_name
		FROM employees e
		LEFT JOIN job_positions jp ON e.job_position_id = jp.id
		WHERE e.deleted_at IS NULL`

	if status != "" {
		query += fmt.Sprintf(" AND e.status = '%s'", status)
	}

	if err := db.DB.Raw(query).Scan(&results).Error; err != nil {
		log.Printf("[EmployeeService] Export query error: %v", err)
		return
	}

	var fileData []byte
	var err error
	ext := "csv"

	if strings.ToLower(format) == "pdf" {
		ext = "pdf"
		fileData, err = s.generateEmployeePDF(results, status)
	} else {
		buf := new(bytes.Buffer)
		writer := csv.NewWriter(buf)
		writer.Write([]string{"Employee NO", "ID NO", "Employee Names", "Dob", "Gender", "Marital Status", "Position"})

		for _, e := range results {
			dob := ""
			if e.DateOfBirth != nil {
				dob = e.DateOfBirth.Format("2006-01-02")
			}
			writer.Write([]string{e.EmployeeNo, e.IdNo, e.FullName, dob, e.Gender, e.MaritalStatus, e.PositionName})
		}
		writer.Flush()
		fileData = buf.Bytes()
		err = writer.Error()
	}

	if err != nil {
		log.Printf("[EmployeeService] Error generating export: %v", err)
		return
	}

	exportDir := "./storage/exports"
	os.MkdirAll(exportDir, 0755)
	filename := fmt.Sprintf("employees_export_%d.%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(exportDir, filename)
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		log.Printf("[EmployeeService] Error saving export file: %v", err)
		return
	}

	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            fmt.Sprintf("Employee Export (%s) Ready", strings.ToUpper(ext)),
		Message:          fmt.Sprintf("Your employee data %s export is ready for download.", ext),
		NotificationType: "SUCCESS",
		DownloadLink:     fmt.Sprintf("/api/employees/export/download/%s", filename),
	})
}

func (s *EmployeeService) generateEmployeePDF(results []employeeExportQueryResult, status string) ([]byte, error) {
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
	pdf.CellFormat(0, 10, "EMPLOYEE REGISTER", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 8)
	headers := []string{"Emp-No", "ID No", "Full Name", "DOB", "Gender", "Marital Status", "Position"}
	widths := []float64{25, 25, 70, 25, 20, 30, 45}

	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 8)
	for _, e := range results {
		dob := ""
		if e.DateOfBirth != nil {
			dob = e.DateOfBirth.Format("2006-01-02")
		}
		pdf.CellFormat(widths[0], 8, e.EmployeeNo, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 8, e.IdNo, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 8, e.FullName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 8, dob, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[4], 8, e.Gender, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[5], 8, e.MaritalStatus, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[6], 8, e.PositionName, "1", 0, "L", false, 0, "")
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}
