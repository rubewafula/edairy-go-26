package services

import (
	"log"
	"math"
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
)

const (
	MilkCollectionOverviewDefaultDays = 7
	MilkCollectionOverviewMaxDays     = 90
)

type MilkCollectionOverviewService struct{}

func NewMilkCollectionOverviewService() *MilkCollectionOverviewService {
	return &MilkCollectionOverviewService{}
}

// collectionDateExpr is the business date of a member collection: the journal date,
// falling back to when the entry was recorded.
const collectionDateExpr = "DATE(COALESCE(mj.journal_date, mje.created_at))"

type dailyQuantityRow struct {
	Date     string
	Quantity float64
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

// MilkOverviewScope limits the overview to one member or one transporter. Kind is
// "all" (staff: whole cooperative), "member" or "transporter".
type MilkOverviewScope struct {
	Kind string
	ID   uint64 // member_registrations.id or transporters.id; 0 = not linked (shows nothing)
	Name string
}

// ResolveMilkOverviewScope works out whose milk the logged-in user may see.
// A user with a Member role sees only their own deliveries; a Transporter role sees only
// the milk on their journals; everyone else (staff) sees the whole cooperative.
// Members and transporters are linked to their login by email (member_registrations.email,
// transporters.email_address), falling back to the phone on an individual transporter's member.
func ResolveMilkOverviewScope(userID uint64) MilkOverviewScope {
	var user struct {
		ID    uint64
		Email string
	}
	db.DB.Raw(`SELECT id, email FROM users WHERE id = ? AND deleted_at IS NULL`, userID).Scan(&user)
	var roles []string
	db.DB.Raw(`SELECT LOWER(r.name) FROM roles r JOIN user_roles ur ON ur.role_id = r.id WHERE ur.user_id = ?`, userID).Scan(&roles)
	has := func(word string) bool {
		for _, r := range roles {
			if strings.Contains(r, word) {
				return true
			}
		}
		return false
	}
	email := strings.TrimSpace(user.Email)
	// Farmer logins (POST /api/member-login) are linked by users.member_id.
	if memberID := MemberIDForUser(userID); memberID > 0 {
		var name string
		db.DB.Raw(`SELECT TRIM(CONCAT_WS(' ', first_name, last_name)) FROM member_registrations WHERE id = ?`, memberID).Scan(&name)
		return MilkOverviewScope{Kind: "member", ID: memberID, Name: name}
	}
	staff := has("admin") || has("manager") || has("clerk") || has("accountant") || has("employee") || has("staff") || has("grader")
	if has("transporter") && !staff {
		var t struct {
			ID   uint64
			Name string
		}
		db.DB.Raw(`SELECT t.id, t.transporter_no AS name FROM transporters t
			WHERE t.deleted_at IS NULL AND LOWER(t.email_address) = LOWER(?) AND ? <> ''
			ORDER BY t.id LIMIT 1`, email, email).Scan(&t)
		if t.ID == 0 {
			log.Printf("[MilkOverview] user %d has a transporter role but no transporter with email %q", userID, email)
		}
		return MilkOverviewScope{Kind: "transporter", ID: t.ID, Name: t.Name}
	}
	if has("member") && !staff {
		var m struct {
			ID   uint64
			Name string
		}
		db.DB.Raw(`SELECT id, TRIM(CONCAT_WS(' ', first_name, last_name)) AS name FROM member_registrations
			WHERE deleted_at IS NULL AND LOWER(email) = LOWER(?) AND ? <> ''
			ORDER BY id LIMIT 1`, email, email).Scan(&m)
		if m.ID == 0 {
			log.Printf("[MilkOverview] user %d has a member role but no member with email %q", userID, email)
		}
		return MilkOverviewScope{Kind: "member", ID: m.ID, Name: m.Name}
	}
	return MilkOverviewScope{Kind: "all"}
}

// GetOverview returns milk collections for the last `days` days (today inclusive),
// compared with the `days` days before that, limited to the given scope.
func (s *MilkCollectionOverviewService) GetOverview(days int, scope MilkOverviewScope) dtos.MilkCollectionOverviewResponse {
	// Extra WHERE for the journal-entry queries and for rejects.
	entryFilter, rejectFilter := "", ""
	var scopeArgs []interface{}
	switch scope.Kind {
	case "member":
		entryFilter, rejectFilter = " AND mje.member_id = ?", " AND member_id = ?"
		scopeArgs = []interface{}{scope.ID}
	case "transporter":
		entryFilter, rejectFilter = " AND mj.transporter_id = ?", " AND transporter_id = ?"
		scopeArgs = []interface{}{scope.ID}
	}
	if days <= 0 {
		days = MilkCollectionOverviewDefaultDays
	}
	if days > MilkCollectionOverviewMaxDays {
		days = MilkCollectionOverviewMaxDays
	}

	const layout = "2006-01-02"
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start := today.AddDate(0, 0, -(days - 1))
	prevStart := start.AddDate(0, 0, -days)
	endExclusive := today.AddDate(0, 0, 1)

	var rows []dailyQuantityRow
	db.DB.Raw(`
		SELECT
			DATE_FORMAT(`+collectionDateExpr+`, '%Y-%m-%d') AS date,
			COALESCE(SUM(IFNULL(mje.quantity, 0)), 0) AS quantity
		FROM milk_journal_entries mje
		LEFT JOIN milk_journals mj ON mj.id = mje.milk_journal_id AND mj.deleted_at IS NULL
		WHERE mje.deleted_at IS NULL
			AND `+collectionDateExpr+` >= ?
			AND `+collectionDateExpr+` < ?`+entryFilter+`
		GROUP BY `+collectionDateExpr+`
	`, append([]interface{}{prevStart.Format(layout), endExclusive.Format(layout)}, scopeArgs...)...).Scan(&rows)

	byDate := make(map[string]float64, len(rows))
	for _, r := range rows {
		byDate[r.Date] = r.Quantity
	}

	trend := make([]dtos.MilkCollectionOverviewPoint, 0, days)
	var total, previous float64
	for d := start; d.Before(endExclusive); d = d.AddDate(0, 0, 1) {
		key := d.Format(layout)
		qty := byDate[key]
		total += qty
		trend = append(trend, dtos.MilkCollectionOverviewPoint{Date: key, Quantity: round2(qty)})
	}
	for d := prevStart; d.Before(start); d = d.AddDate(0, 0, 1) {
		previous += byDate[d.Format(layout)]
	}

	var changePercent *float64
	if previous > 0 {
		v := round2((total - previous) / previous * 100)
		changePercent = &v
	}

	var activeMembers int64
	db.DB.Raw(`
		SELECT COUNT(DISTINCT mje.member_id)
		FROM milk_journal_entries mje
		LEFT JOIN milk_journals mj ON mj.id = mje.milk_journal_id AND mj.deleted_at IS NULL
		WHERE mje.deleted_at IS NULL
			AND `+collectionDateExpr+` >= ?
			AND `+collectionDateExpr+` < ?`+entryFilter+`
	`, append([]interface{}{start.Format(layout), endExclusive.Format(layout)}, scopeArgs...)...).Scan(&activeMembers)

	var rejected float64
	db.DB.Raw(`
		SELECT COALESCE(SUM(IFNULL(quantity, 0)), 0)
		FROM milk_rejects
		WHERE deleted_at IS NULL
			AND DATE(COALESCE(transaction_date, created_at)) >= ?
			AND DATE(COALESCE(transaction_date, created_at)) < ?`+rejectFilter+`
	`, append([]interface{}{start.Format(layout), endExclusive.Format(layout)}, scopeArgs...)...).Scan(&rejected)

	return dtos.MilkCollectionOverviewResponse{
		PeriodDays:            days,
		StartDate:             start.Format(layout),
		EndDate:               today.Format(layout),
		TotalQuantity:         round2(total),
		DailyAverage:          round2(total / float64(days)),
		TodayQuantity:         round2(byDate[today.Format(layout)]),
		PreviousTotalQuantity: round2(previous),
		ChangePercent:         changePercent,
		RejectedQuantity:      round2(rejected),
		ActiveMembers:         activeMembers,
		Trend:                 trend,
		Scope:                 scope.Kind,
		ScopeID:               scope.ID,
		ScopeName:             scope.Name,
	}
}

// EnsureBankBranchViewPermission lets every role that can list banks also list bank branches
// (the mobile app downloads both as common data for member registration).
func EnsureBankBranchViewPermission() {
	db.DB.Exec(`INSERT INTO permissions (name, guard_name, created_at, updated_at)
		SELECT 'bank-branches.view', 'web', NOW(), NOW() FROM DUAL
		WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE name = 'bank-branches.view')`)
	if err := db.DB.Exec(`INSERT INTO role_permissions (role_id, permission_id)
		SELECT DISTINCT rp.role_id, bp.id FROM role_permissions rp
		JOIN permissions src ON src.id = rp.permission_id AND src.name = 'banks.view' AND src.deleted_at IS NULL
		JOIN permissions bp ON bp.name = 'bank-branches.view' AND bp.deleted_at IS NULL
		WHERE NOT EXISTS (SELECT 1 FROM role_permissions x WHERE x.role_id = rp.role_id AND x.permission_id = bp.id)`).Error; err != nil {
		log.Printf("[Banks] grant bank-branches.view: %v", err)
	}
}

// EnsureMilkOverviewPermission grants milk-collection-overview.view to every role. The endpoint
// scopes itself (a member sees only their milk, a transporter only theirs), so it is safe for all.
func EnsureMilkOverviewPermission() {
	const name = "milk-collection-overview.view"
	if err := db.DB.Exec(`INSERT INTO permissions (name, guard_name, created_at, updated_at)
		SELECT ?, 'web', NOW(), NOW() FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE name = ?)`, name, name).Error; err != nil {
		log.Printf("[MilkOverview] seed permission: %v", err)
		return
	}
	if err := db.DB.Exec(`INSERT INTO role_permissions (role_id, permission_id)
		SELECT r.id, p.id FROM roles r JOIN permissions p ON p.name = ? AND p.deleted_at IS NULL
		WHERE r.deleted_at IS NULL
		AND NOT EXISTS (SELECT 1 FROM role_permissions rp WHERE rp.role_id = r.id AND rp.permission_id = p.id)`, name).Error; err != nil {
		log.Printf("[MilkOverview] grant permission: %v", err)
	}
}
