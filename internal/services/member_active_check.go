package services

import (
	"log"
	"net/http"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

// requireActiveMember refuses loans (applications, approvals) for members who are not ACTIVE.
func requireActiveMember(memberID uint64) error {
	var m models.Member
	if err := db.DB.Select("id, status").First(&m, memberID).Error; err != nil {
		return serviceErr(http.StatusNotFound, "member not found")
	}
	if !memberCanBorrow(m.Status) {
		return serviceErr(http.StatusForbidden, "%s", OnlyActiveMembersMessage)
	}
	return nil
}

// refreshMemberBorrowStatus re-checks one member's can_borrow (e.g. after a status change).
func refreshMemberBorrowStatus(memberID uint64) {
	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		return NewMemberCreditService().RefreshBorrowStatus(tx, memberID)
	}); err != nil {
		log.Printf("[MemberCredit] refresh borrow status for member %d: %v", memberID, err)
	}
}

// EnsureOnlyActiveMembersBorrow re-checks every member whose stored can_borrow disagrees with the
// active-only rule (run at startup so existing pending/inactive members stop being offered loans).
func EnsureOnlyActiveMembersBorrow() {
	var ids []uint64
	db.DB.Raw(`SELECT id FROM member_registrations WHERE deleted_at IS NULL
		AND ((can_borrow = 1 AND UPPER(TRIM(COALESCE(status, ''))) <> 'ACTIVE')
		  OR (UPPER(TRIM(COALESCE(status, ''))) = 'ACTIVE' AND cannot_borrow_code = 'MEMBER_STATUS'))`).Scan(&ids)
	for _, id := range ids {
		refreshMemberBorrowStatus(id)
	}
	if len(ids) > 0 {
		log.Printf("[MemberCredit] borrow status re-checked for %d member(s): only active members may borrow", len(ids))
	}
}
