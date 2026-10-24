-- Set every member's credit limit to KES 5,000 (requested 2026-10-18).
-- Run against the `edairy` database. Run the updated API once first so the
-- can_borrow / cannot_borrow_code / cannot_borrow_reason columns exist.
--
-- Undo:  UPDATE member_registrations m JOIN member_credit_limit_backup_20261018 b ON b.id = m.id
--        SET m.credit_limit = b.credit_limit, m.credit_limit_updated_at = b.credit_limit_updated_at;

START TRANSACTION;

-- 1. Backup of the current values (kept so this can be undone).
CREATE TABLE IF NOT EXISTS member_credit_limit_backup_20261018 AS
SELECT id, member_no, credit_limit, credit_limit_updated_at, NOW() AS backed_up_at
FROM member_registrations;

-- 2. Audit trail, same table the API writes to.
INSERT INTO credit_limit_change_log
    (customer_id, customer_type, old_limit, credit_limit, amount, action, reference,
     source_type, source_id, narration, created_by, updated_by, created_at, updated_at)
SELECT id, 'member', COALESCE(credit_limit, 0), 5000, 5000 - COALESCE(credit_limit, 0), 'SEEDED',
       'BULK-SET-5000-20261018', 'manual_bulk_update', 0,
       'All member credit limits set to 5000', 0, 0, NOW(), NOW()
FROM member_registrations
WHERE deleted_at IS NULL AND COALESCE(credit_limit, 0) <> 5000;

-- 3. The change itself.
UPDATE member_registrations
SET credit_limit = 5000, credit_limit_updated_at = NOW()
WHERE deleted_at IS NULL;

-- 4. Refresh can_borrow with the API's rules (status → last application → limit vs minimum).
SET @min_borrow := COALESCE(
    (SELECT min_amount FROM loan_products WHERE is_default_member_loan = 1 AND is_active = 1 ORDER BY id LIMIT 1),
    (SELECT min_amount FROM loan_products WHERE is_active = 1 ORDER BY id LIMIT 1),
    0);

UPDATE member_registrations m
LEFT JOIN (
    SELECT la.member_id, la.id, la.status, la.source
    FROM loan_applications la
    JOIN (SELECT member_id, MAX(id) AS last_id FROM loan_applications GROUP BY member_id) x
      ON x.last_id = la.id
) last_app ON last_app.member_id = m.id
SET
  m.cannot_borrow_code = CASE
    WHEN UPPER(TRIM(COALESCE(m.status, ''))) IN ('SUSPENDED','INACTIVE','EXITED','DECEASED','DELETED','CLOSED','REJECTED') THEN 'MEMBER_STATUS'
    WHEN last_app.status IN ('SUBMITTED','UNDER_REVIEW')
      OR (last_app.status = 'DRAFT' AND (last_app.source <> 'MOBILE' OR EXISTS (
            SELECT 1 FROM loan_application_otps o
            WHERE o.loan_application_id = last_app.id AND o.status = 'PENDING' AND o.expires_at > NOW())))
      THEN 'PENDING_APPLICATION'
    WHEN 5000 < @min_borrow THEN 'LIMIT_INSUFFICIENT'
    ELSE NULL END,
  m.cannot_borrow_reason = CASE
    WHEN UPPER(TRIM(COALESCE(m.status, ''))) IN ('SUSPENDED','INACTIVE','EXITED','DECEASED','DELETED','CLOSED','REJECTED')
      THEN CONCAT('Member account is ', LOWER(TRIM(m.status)), ' and cannot borrow.')
    WHEN last_app.status IN ('SUBMITTED','UNDER_REVIEW')
      THEN CONCAT('Loan application #', last_app.id, ' is still being processed (', LOWER(REPLACE(last_app.status, '_', ' ')), '). A new cashout can be requested once it is approved, disbursed or rejected.')
    WHEN last_app.status = 'DRAFT' AND (last_app.source <> 'MOBILE' OR EXISTS (
            SELECT 1 FROM loan_application_otps o
            WHERE o.loan_application_id = last_app.id AND o.status = 'PENDING' AND o.expires_at > NOW()))
      THEN CONCAT('Loan application #', last_app.id, ' is awaiting OTP confirmation. Confirm it, or wait for it to expire, before requesting another cashout.')
    WHEN 5000 < @min_borrow
      THEN CONCAT('Credit limit not sufficient: KES 5,000.00 is below the minimum borrowable amount of KES ', FORMAT(@min_borrow, 2), '.')
    ELSE NULL END
WHERE m.deleted_at IS NULL;

UPDATE member_registrations
SET can_borrow = (cannot_borrow_code IS NULL)
WHERE deleted_at IS NULL;

COMMIT;

-- Check:
SELECT credit_limit, can_borrow, cannot_borrow_code, COUNT(*) AS members
FROM member_registrations WHERE deleted_at IS NULL
GROUP BY credit_limit, can_borrow, cannot_borrow_code;
