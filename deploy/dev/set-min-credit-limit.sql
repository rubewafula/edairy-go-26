-- Set every member's credit limit to at least KES 5,000 directly in the DB (dev/test).
--   mysql -uroot -p edairy < deploy/dev/set-min-credit-limit.sql
-- NOTE: the API recomputes credit_limit from unpaid milk on every credit query / milk / loan event,
-- so these values will be overwritten for a member the next time their limit is recalculated.
-- Use `go run ./cmd/seed-credit-limits` for limits that survive recalculation.

SET @min_limit := 5000;
SET @user_id   := 1;

START TRANSACTION;

-- 1. Log the change first (captures the old limit). credit_limit = the new current limit.
INSERT INTO credit_limit_change_log
  (customer_id, customer_type, old_limit, credit_limit, amount, action,
   reference, source_type, source_id, narration, created_by, updated_by, created_at, updated_at)
SELECT m.id, 'member',
       COALESCE(m.credit_limit, 0),
       @min_limit,
       @min_limit - COALESCE(m.credit_limit, 0),
       'SEEDED',
       'MIN-CREDIT-5000', 'manual_seed', 0,
       CONCAT('Credit limit set to minimum KES ', @min_limit, ' at DB level'),
       @user_id, @user_id, NOW(), DATE_FORMAT(NOW(), '%Y-%m-%d %H:%i:%s')
FROM member_registrations m
WHERE m.deleted_at IS NULL
  AND COALESCE(m.credit_limit, 0) < @min_limit;

-- 2. Raise the current limit; clear "no/insufficient limit" borrow blocks (other blocks are kept)
UPDATE member_registrations m
SET m.credit_limit = @min_limit,
    m.credit_limit_updated_at = NOW(),
    m.can_borrow = CASE WHEN m.cannot_borrow_code IN ('NO_CREDIT_LIMIT', 'LIMIT_INSUFFICIENT') OR m.cannot_borrow_code IS NULL
                        THEN 1 ELSE m.can_borrow END,
    m.cannot_borrow_reason = CASE WHEN m.cannot_borrow_code IN ('NO_CREDIT_LIMIT', 'LIMIT_INSUFFICIENT')
                                  THEN NULL ELSE m.cannot_borrow_reason END,
    m.cannot_borrow_code = CASE WHEN m.cannot_borrow_code IN ('NO_CREDIT_LIMIT', 'LIMIT_INSUFFICIENT')
                                THEN NULL ELSE m.cannot_borrow_code END
WHERE m.deleted_at IS NULL
  AND COALESCE(m.credit_limit, 0) < @min_limit;

COMMIT;

-- 3. Check: should return 0
SELECT COUNT(*) AS members_below_min
FROM member_registrations
WHERE deleted_at IS NULL AND COALESCE(credit_limit, 0) < @min_limit;
