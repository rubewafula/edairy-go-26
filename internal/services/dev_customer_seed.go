package services

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rubewafula/edairy-go-26/internal/db"
)

// devTestCustomers are milk buyers for testing Customer Milk Delivery. Seeded only when
// ENV=development (set SEED_TEST_CUSTOMERS=false to skip). Safe to re-run: rows are keyed by customer_no.
var devTestCustomers = []struct{ No, Name, Phone, Type, Town string; Rate, Limit float64; Terms int }{
	{"DEV-CUS-001", "Kiambu Fresh Dairies Ltd", "254711000001", "Company", "Kiambu", 52, 500000, 30},
	{"DEV-CUS-002", "Githunguri Milk Bar", "254711000002", "Shop", "Githunguri", 55, 50000, 7},
	{"DEV-CUS-003", "Limuru Hotel & Lodge", "254711000003", "Hotel", "Limuru", 58, 80000, 14},
	{"DEV-CUS-004", "St. Mary's Boarding School", "254711000004", "School", "Thika", 54, 120000, 30},
	{"DEV-CUS-005", "Mama Wanjiru Kiosk", "254711000005", "Shop", "Ruiru", 60, 10000, 0},
	{"DEV-CUS-006", "Nyeri Yoghurt Makers", "254711000006", "Company", "Nyeri", 50, 300000, 21},
	{"DEV-CUS-007", "Juja Supermarket", "254711000007", "Shop", "Juja", 56, 150000, 14},
	{"DEV-CUS-008", "Karatina Cafe", "254711000008", "Hotel", "Karatina", 59, 20000, 7},
	{"DEV-CUS-009", "Kikuyu Hospital Canteen", "254711000009", "Hospital", "Kikuyu", 55, 90000, 30},
	{"DEV-CUS-010", "Walk-in Cash Buyer", "254711000010", "Individual", "Kiambu", 60, 0, 0},
}

// devCustomerTypeDescriptions are the customer types used by the test customers. Existing types with
// the same name (any case, e.g. "School", "Individual" from earlier data) are reused, not duplicated.
var devCustomerTypeDescriptions = map[string]string{
	"Individual": "Private person buying milk for household use",
	"School":     "Schools and colleges (boarding/day)",
	"Hotel":      "Hotels, restaurants and cafes",
	"Hospital":   "Hospitals and health centres",
	"Shop":       "Shops, kiosks, milk bars and supermarkets reselling milk",
	"Company":    "Processors and companies buying in bulk",
}

// SeedDevCustomers inserts the test customers (and their customer types) if they are missing.
func SeedDevCustomers() {
	if os.Getenv("ENV") != "development" || strings.EqualFold(os.Getenv("SEED_TEST_CUSTOMERS"), "false") {
		return
	}
	added := 0
	for _, c := range devTestCustomers {
		db.DB.Exec(`INSERT INTO customer_types (name, description, created_at, updated_at)
			SELECT ?, ?, NOW(), NOW() FROM DUAL
			WHERE NOT EXISTS (SELECT 1 FROM customer_types WHERE LOWER(name) = LOWER(?) AND deleted_at IS NULL)`,
			c.Type, devCustomerTypeDescriptions[c.Type], c.Type)
		email := strings.ToLower(strings.ReplaceAll(c.No, "-", "")) + "@customers.test"
		res := db.DB.Exec(`INSERT INTO customers (customer_type_id, full_names, phone, email_address, customer_no,
				status, credit_limit, rate, payment_terms_days, postal_town, created_at, updated_at)
			SELECT (SELECT id FROM customer_types WHERE LOWER(name) = LOWER(?) AND deleted_at IS NULL ORDER BY id LIMIT 1),
				?, ?, ?, ?, 'ACTIVE', ?, ?, ?, ?, NOW(), NOW() FROM DUAL
			WHERE NOT EXISTS (SELECT 1 FROM customers WHERE customer_no = ?)`,
			c.Type, c.Name, c.Phone, email, c.No, c.Limit, c.Rate, c.Terms, c.Town, c.No)
		if res.Error != nil {
			log.Printf("[DevSeed] customer %s: %v", c.No, res.Error)
			continue
		}
		added += int(res.RowsAffected)
		// Keep already-seeded test customers on the right type.
		db.DB.Exec(`UPDATE customers SET customer_type_id =
				(SELECT id FROM customer_types WHERE LOWER(name) = LOWER(?) AND deleted_at IS NULL ORDER BY id LIMIT 1)
			WHERE customer_no = ?`, c.Type, c.No)
	}
	// Remove the placeholder types an earlier seed created, once nothing uses them.
	db.DB.Exec(`UPDATE customer_types ct SET ct.deleted_at = NOW()
		WHERE ct.deleted_at IS NULL AND ct.description LIKE '% (test data)'
		AND NOT EXISTS (SELECT 1 FROM customers c WHERE c.customer_type_id = ct.id AND c.deleted_at IS NULL)`)
	// Undo an earlier soft delete so the test customers always show.
	db.DB.Exec(`UPDATE customers SET deleted_at = NULL, status = 'ACTIVE' WHERE customer_no LIKE 'DEV-CUS-%'`)
	log.Printf("[DevSeed] test customers: %s added, %d in the seed set", fmt.Sprint(added), len(devTestCustomers))
}
