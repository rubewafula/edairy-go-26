package middleware

import (
	"encoding/json"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type ActivityLoggerPlugin struct{}

func NewActivityLoggerPlugin() *ActivityLoggerPlugin {
	return &ActivityLoggerPlugin{}
}

func (p *ActivityLoggerPlugin) Name() string {
	return "activity-logger-plugin"
}

func (p *ActivityLoggerPlugin) Initialize(db *gorm.DB) error {

	db.Callback().
		Create().
		After("gorm:create").
		Register("activity_logger_create", func(tx *gorm.DB) {
			p.writeLog(tx, "CREATE")
		})

	db.Callback().
		Update().
		After("gorm:update").
		Register("activity_logger_update", func(tx *gorm.DB) {
			p.writeLog(tx, "UPDATE")
		})

	db.Callback().
		Delete().
		After("gorm:delete").
		Register("activity_logger_delete", func(tx *gorm.DB) {
			p.writeLog(tx, "DELETE")
		})

	return nil
}

func (p *ActivityLoggerPlugin) writeLog(tx *gorm.DB, operation string) {

	// -----------------------------------
	// Prevent recursive logging
	// -----------------------------------

	if tx.Statement.Table == "activity_log" {
		return
	}

	// -----------------------------------
	// Skip failed queries
	// -----------------------------------

	if tx.Error != nil {
		return
	}

	// -----------------------------------
	// Table name
	// -----------------------------------

	tableName := tx.Statement.Table

	// -----------------------------------
	// Authenticated user
	// -----------------------------------

	var causerID *uint64
	var causerType *string

	if tx.Statement.Context != nil {

		if authUser, ok := tx.Statement.Context.Value(UserContextKey).(AuthUser); ok {

			causerID = &authUser.UserID

			userType := "User"
			causerType = &userType
		}
	}

	// -----------------------------------
	// Build properties JSON
	// -----------------------------------

	properties, _ := json.Marshal(map[string]interface{}{
		"table":         tableName,
		"operation":     operation,
		"rows_affected": tx.RowsAffected,
		"sql":           tx.Statement.SQL.String(),
		"vars":          tx.Statement.Vars,
	})

	// -----------------------------------
	// Build log entry
	// -----------------------------------

	now := time.Now()

	logName := "database"
	description := operation + " on " + tableName
	event := operation

	logEntry := models.ActivityLog{
		LogName:     &logName,
		Description: description,
		SubjectType: &tableName,
		CauserType:  causerType,
		CauserID:    causerID,
		Properties:  properties,
		Event:       &event,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	// -----------------------------------
	// Save activity log
	// -----------------------------------

	tx.Session(&gorm.Session{
		NewDB:     true,
		SkipHooks: true,
	}).Table("activity_log").Create(&logEntry)
}
