-- TEST ONLY: let every member borrow (requested 2026-10-19).
-- Run against the `edairy` database with the updated API (can_borrow columns exist).
-- Note: the API recomputes can_borrow for a member on their next milk, payroll or loan event.

START TRANSACTION;

-- Backup so this can be undone.
CREATE TABLE IF NOT EXISTS member_can_borrow_backup_20261019 AS
SELECT id, member_no, credit_limit, can_borrow, cannot_borrow_code, cannot_borrow_reason, NOW() AS backed_up_at
FROM member_registrations;

UPDATE member_registrations
SET can_borrow = 1,
    cannot_borrow_code = NULL,
    cannot_borrow_reason = NULL
WHERE deleted_at IS NULL;

-- Loan requests are also checked against the stored credit limit (total repayable must not exceed it).
-- Members below 5000 are raised to 5000 so a test loan isn't rejected for limit.
UPDATE member_registrations
SET credit_limit = 5000, credit_limit_updated_at = NOW()
WHERE deleted_at IS NULL AND COALESCE(credit_limit, 0) < 5000;

COMMIT;

-- Check:
SELECT can_borrow, credit_limit >= 5000 AS limit_ok, COUNT(*) AS members
FROM member_registrations WHERE deleted_at IS NULL GROUP BY 1, 2;

-- Members whose LAST loan application is still unprocessed. A new loan request is still
-- refused for these until that application is approved, disbursed, rejected or cancelled.
SELECT la.member_id, la.id AS application_id, la.status, la.requested_amount
FROM loan_applications la
JOIN (SELECT member_id, MAX(id) AS last_id FROM loan_applications GROUP BY member_id) x ON x.last_id = la.id
WHERE la.status IN ('DRAFT', 'SUBMITTED', 'UNDER_REVIEW');

-- Undo:
-- UPDATE member_registrations m JOIN member_can_borrow_backup_20261019 b ON b.id = m.id
-- SET m.can_borrow = b.can_borrow, m.cannot_borrow_code = b.cannot_borrow_code,
--     m.cannot_borrow_reason = b.cannot_borrow_reason, m.credit_limit = b.credit_limit;
