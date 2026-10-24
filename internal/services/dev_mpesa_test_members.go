package services

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/rubewafula/edairy-go-26/internal/db"
)

// devMpesaTestMembers are the members used for real M-Pesa testing in development. Their phones are
// the payout allow-list (see dev_payout_allowlist.go), so loans taken by them can actually be paid.
// Seeded only when ENV=development (SEED_TEST_MEMBERS=false skips). Safe to re-run: keyed by member_no.
var devMpesaTestMembers = []struct{ No, IDNo, First, Last, Phone string }{
	{"DEV-MEM-001", "DEVMPESA001", "Mpesa", "Tester One", "254705182016"},
	{"DEV-MEM-002", "DEVMPESA002", "Mpesa", "Tester Two", "254726986944"},
	{"DEV-MEM-003", "DEVMPESA003", "Mpesa", "Tester Three", "254721695524"},
}

// SeedDevMpesaTestMembers inserts the M-Pesa test members and, on the local "edairy" database only,
// tops their credit limit up to SEED_TEST_MEMBER_CREDIT (default KES 10,000) with seeded unpaid milk.
func SeedDevMpesaTestMembers() {
	if os.Getenv("ENV") != "development" || strings.EqualFold(os.Getenv("SEED_TEST_MEMBERS"), "false") {
		return
	}
	// Put them on the same route/member type as a member who can already borrow (advance rate set).
	var ref struct{ RouteID, MemberTypeID uint64 }
	db.DB.Raw(`SELECT COALESCE(route_id,0) AS route_id, COALESCE(member_type_id,0) AS member_type_id
		FROM member_registrations WHERE deleted_at IS NULL AND can_borrow = 1 AND route_id > 0
		ORDER BY id LIMIT 1`).Scan(&ref)
	if ref.RouteID == 0 {
		db.DB.Raw(`SELECT id FROM routes WHERE deleted_at IS NULL ORDER BY id LIMIT 1`).Scan(&ref.RouteID)
	}
	if ref.MemberTypeID == 0 {
		db.DB.Raw(`SELECT id FROM member_types WHERE deleted_at IS NULL ORDER BY id LIMIT 1`).Scan(&ref.MemberTypeID)
	}

	added := 0
	var ids []uint64
	for _, m := range devMpesaTestMembers {
		res := db.DB.Exec(`INSERT INTO member_registrations (member_no, id_no, first_name, last_name, gender,
				primary_phone, status, route_id, member_type_id, date_registered, created_at, updated_at)
			SELECT ?, ?, ?, ?, 'Male', ?, 'active', NULLIF(?,0), NULLIF(?,0), NOW(), NOW(), NOW() FROM DUAL
			WHERE NOT EXISTS (SELECT 1 FROM member_registrations WHERE member_no = ?)`,
			m.No, m.IDNo, m.First, m.Last, m.Phone, ref.RouteID, ref.MemberTypeID, m.No)
		if res.Error != nil {
			log.Printf("[DevSeed] test member %s: %v", m.No, res.Error)
			continue
		}
		added += int(res.RowsAffected)
		// Keep them usable: undelete, active, correct phone.
		db.DB.Exec(`UPDATE member_registrations SET deleted_at = NULL, status = 'active', primary_phone = ?
			WHERE member_no = ?`, m.Phone, m.No)
		var id uint64
		db.DB.Raw(`SELECT id FROM member_registrations WHERE member_no = ?`, m.No).Scan(&id)
		if id > 0 {
			ids = append(ids, id)
		}
	}
	log.Printf("[DevSeed] M-Pesa test members: %d added, ids %v (phones %s)", added, ids, strings.Join(DevPayoutPhones(), ", "))

	dbName, dbHost := os.Getenv("DB_NAME"), os.Getenv("DB_HOST")
	if len(ids) == 0 || dbName != "edairy" || (dbHost != "localhost" && dbHost != "127.0.0.1") {
		return // seeded milk is payable: only on the local dev database
	}
	minLimit := 10000.0
	if v, err := strconv.ParseFloat(os.Getenv("SEED_TEST_MEMBER_CREDIT"), 64); err == nil && v > 0 {
		minLimit = v
	}
	res, err := SeedMinimumCreditLimitsFor(ids, minLimit, 50, 1)
	if err != nil {
		log.Printf("[DevSeed] M-Pesa test members credit: %v", err)
		return
	}
	log.Printf("[DevSeed] M-Pesa test members credit >= KES %s: raised %d, already above %d%s",
		formatKES(minLimit), res.Raised, res.AlreadyAbove, failuresSuffix(res.Failures))
}

func failuresSuffix(f []string) string {
	if len(f) == 0 {
		return ""
	}
	return fmt.Sprintf("; failed: %s", strings.Join(f, "; "))
}
