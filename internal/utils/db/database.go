package db

import appdb "github.com/rubewafula/edairy-go-26/internal/db"

// Deprecated: use github.com/rubewafula/edairy-go-26/internal/db instead.

func ConnectToDatabase() {
	appdb.ConnectToDatabase()
}

func CloseDatabase() error {
	return appdb.CloseDatabase()
}
