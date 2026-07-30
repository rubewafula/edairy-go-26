-- COA cleanup and posting rule fixes for dairy deployments (P0)
-- Review before running on production; backs up affected rows first.

-- 1. Deactivate misclassified duplicate accounts (444 INCOME type cash accounts)
UPDATE accounts a
JOIN account_categories ac ON ac.id = a.account_category_id
JOIN account_types at ON at.id = ac.account_type_id
SET a.is_active = 0, a.is_postable = 0
WHERE a.account_code = '444' AND at.name = 'INCOME';

-- 2. Deduplicate posting rules — keep lowest id per transaction_type
DELETE tpr FROM transaction_posting_rules tpr
INNER JOIN (
  SELECT transaction_type, MIN(id) AS keep_id
  FROM transaction_posting_rules
  GROUP BY transaction_type
  HAVING COUNT(*) > 1
) d ON d.transaction_type = tpr.transaction_type AND tpr.id <> d.keep_id;

-- 3. Fix STORE_SALE_REVENUE to debit Cash at Bank (1000) instead of 444
UPDATE transaction_posting_rules tpr
JOIN accounts bad ON bad.id = tpr.debit_account_id AND bad.account_code = '444'
JOIN accounts good ON good.account_code = '1000' AND good.deleted_at IS NULL
SET tpr.debit_account_id = good.id
WHERE tpr.transaction_type = 'STORE_SALE_REVENUE';

-- 4. Fix remittance rules crediting 444 → 1000
UPDATE transaction_posting_rules tpr
JOIN accounts bad ON bad.id = tpr.credit_account_id AND bad.account_code = '444'
JOIN accounts good ON good.account_code = '1000' AND good.deleted_at IS NULL
SET tpr.credit_account_id = good.id
WHERE tpr.transaction_type IN ('EMPLOYEE_NET_SALARY_PAYMENT','NHIF_REMITTANCE','NSSF_REMITTANCE','PAYE_REMITTANCE');

-- 5. Add missing Go posting rules (requires accounts to exist)
INSERT IGNORE INTO transaction_posting_rules (transaction_type, debit_account_id, credit_account_id, description, created_at, updated_at)
SELECT 'TRANSPORTER_PAYMENT_GROSS', e.id, s.id, 'Transporter gross payment accrual', NOW(), NOW()
FROM accounts e JOIN accounts s ON s.account_code = '2210'
WHERE e.account_code = '5056'
  AND NOT EXISTS (SELECT 1 FROM transaction_posting_rules WHERE transaction_type = 'TRANSPORTER_PAYMENT_GROSS');

INSERT IGNORE INTO transaction_posting_rules (transaction_type, debit_account_id, credit_account_id, description, created_at, updated_at)
SELECT 'TRANSPORTER_BENEFIT_EXPENSE', e.id, s.id, 'Transporter benefit expense', NOW(), NOW()
FROM accounts e JOIN accounts s ON s.account_code = '2215'
WHERE e.account_code = '5100'
  AND NOT EXISTS (SELECT 1 FROM transaction_posting_rules WHERE transaction_type = 'TRANSPORTER_BENEFIT_EXPENSE');

INSERT IGNORE INTO transaction_posting_rules (transaction_type, debit_account_id, credit_account_id, description, created_at, updated_at)
SELECT 'SHARES_IMPORT', b.id, sh.id, 'Opening share balance import', NOW(), NOW()
FROM accounts b JOIN accounts sh ON sh.account_code = '3000'
WHERE b.account_code = '1000'
  AND NOT EXISTS (SELECT 1 FROM transaction_posting_rules WHERE transaction_type = 'SHARES_IMPORT');

INSERT IGNORE INTO transaction_posting_rules (transaction_type, debit_account_id, credit_account_id, description, created_at, updated_at)
SELECT 'STOCK_IMPORT', inv.id, b.id, 'Inventory opening import', NOW(), NOW()
FROM accounts inv JOIN accounts b ON b.account_code = '1000'
WHERE inv.account_code = '1500'
  AND NOT EXISTS (SELECT 1 FROM transaction_posting_rules WHERE transaction_type = 'STOCK_IMPORT');

-- 6. Recommended new accounts (skip if code exists)
INSERT IGNORE INTO accounts (account_code, name, description, account_category_id, is_postable, is_active, created_at, updated_at)
SELECT '1050', 'Petty Cash', 'Cash on hand', ac.id, 1, 1, NOW(), NOW()
FROM account_categories ac JOIN account_types at ON at.id = ac.account_type_id
WHERE at.name = 'ASSET' AND ac.name LIKE '%Cash%' LIMIT 1;

INSERT IGNORE INTO accounts (account_code, name, description, account_category_id, is_postable, is_active, created_at, updated_at)
SELECT '1600', 'Accumulated Depreciation', 'Fixed asset contra', ac.id, 1, 1, NOW(), NOW()
FROM account_categories ac JOIN account_types at ON at.id = ac.account_type_id
WHERE at.name = 'ASSET' AND ac.name LIKE '%Depreciation%' LIMIT 1;
