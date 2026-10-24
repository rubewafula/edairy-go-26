// Command seed-credit-limits gives every member a credit limit of at least KES 5,000 (DEV/TEST ONLY).
//
//	go run ./cmd/seed-credit-limits
//
// Env:
//
//	SEED_MIN_CREDIT_LIMIT  minimum limit per member (default 5000)
//	SEED_MILK_RATE         KES/kg used to create a global default milk rate if none resolves (default 50)
//	SEED_USER_ID           user recorded as created_by in records and credit_limit_change_log (default 1)
//	SEED_ALLOW_DB          set to DB_NAME to allow a database other than the local "edairy"
//
// It seeds confirmed, unpaid milk for the shortfall and recalculates each member, which writes the
// credit_limit_change_log entries. The seeded milk will be paid by the next payroll.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/initializers"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

func envFloat(key string, def float64) float64 {
	if v, err := strconv.ParseFloat(os.Getenv(key), 64); err == nil {
		return v
	}
	return def
}

func main() {
	initializers.LoadEnvVariables()

	dbName, dbHost := os.Getenv("DB_NAME"), os.Getenv("DB_HOST")
	isLocal := dbName == "edairy" && (dbHost == "localhost" || dbHost == "127.0.0.1")
	if !isLocal && os.Getenv("SEED_ALLOW_DB") != dbName {
		fmt.Fprintf(os.Stderr, "refusing to seed %s@%s: this creates payable milk. Set SEED_ALLOW_DB=%s to confirm it is a dev/test database.\n", dbName, dbHost, dbName)
		os.Exit(1)
	}

	db.ConnectToDatabase()

	minLimit := envFloat("SEED_MIN_CREDIT_LIMIT", 5000)
	rate := envFloat("SEED_MILK_RATE", 50)
	userID := uint64(envFloat("SEED_USER_ID", 1))

	res, err := services.SeedMinimumCreditLimits(minLimit, rate, userID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "seed credit limits:", err)
		os.Exit(1)
	}
	fmt.Printf("Members: %d | raised to >= KES %.2f: %d | already above: %d | milk seeded: %.2f kg (journals %s-R<route>)\n",
		res.Members, minLimit, res.Raised, res.AlreadyAbove, res.KilosSeeded, res.JournalPrefix)
	if res.RateCreated {
		fmt.Printf("Created a global default milk rate of KES %.2f/kg (none was configured).\n", rate)
	}
	for _, f := range res.Failures {
		fmt.Println("FAILED:", f)
	}
	if len(res.Failures) > 0 {
		os.Exit(2)
	}
}
