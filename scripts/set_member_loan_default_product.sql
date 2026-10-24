-- Make "Member Loans" the only product members can apply for from the mobile app.
-- The API also does this automatically at startup if exactly one active product is
-- named "Member Loan"/"Member Loans" (or coded MEMBER-LOAN) and none is flagged yet.

-- 1. See the products and which one is flagged:
SELECT id, product_code, product_name, is_active, is_default_member_loan, min_amount, max_amount
FROM loan_products WHERE deleted_at IS NULL ORDER BY id;

-- 2. Flag the member loan product (only one product may carry the flag).
--    If yours has a different name, change the WHERE clause, or use its id: WHERE id = <id>
START TRANSACTION;
UPDATE loan_products SET is_default_member_loan = 0 WHERE is_default_member_loan = 1;
UPDATE loan_products SET is_default_member_loan = 1, is_active = 1
WHERE deleted_at IS NULL AND LOWER(TRIM(product_name)) IN ('member loan', 'member loans')
ORDER BY id LIMIT 1;
COMMIT;

-- 3. Check:
SELECT id, product_name, is_default_member_loan FROM loan_products WHERE is_default_member_loan = 1;
