package dtos

type HRDashboardResponse struct {
	EmployeeStats      HRDashboardEmployeeStats   `json:"employee_stats"`
	LeaveStats         HRDashboardLeaveStats      `json:"leave_stats"`
	EmployeeCategories []HRDashboardCategoryCount `json:"employee_categories"`
}

type HRDashboardEmployeeStats struct {
	TotalEmployees     int64 `json:"total_employees"`
	NewEmployees       int64 `json:"new_employees"`
	ExistingEmployees  int64 `json:"existing_employees"`
	Male               int64 `json:"male"`
	Female             int64 `json:"female"`
}

type HRDashboardLeaveStats struct {
	PendingApplications int64 `json:"pending_applications"`
	EmployeesOnLeave    int64 `json:"employees_on_leave"`
}

type HRDashboardCategoryCount struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}
