package models

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// BeforeCreate gives every transaction header a unique idempotency key. idempotency_key has a unique
// index, so headers created without one (payroll, stores, shares...) all wrote '' and every one after
// the first failed with "Duplicate entry '' for key idx_transactions_idempotency". Such headers get
// AUTO-<reference> (the reference is unique too); ledger postings keep the key they pass in.
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if strings.TrimSpace(t.IdempotencyKey) != "" {
		return nil
	}
	if ref := strings.TrimSpace(t.Reference); ref != "" {
		t.IdempotencyKey = "AUTO-" + ref
	} else {
		t.IdempotencyKey = fmt.Sprintf("AUTO-%d", time.Now().UnixNano())
	}
	return nil
}
