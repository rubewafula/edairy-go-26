-- Full loan-module RBAC permissions (idempotent).
-- Run after 001-010 on each dairy database.

INSERT IGNORE INTO permissions (name, guard_name, created_at, updated_at)
SELECT pname, 'web', NOW(), NOW()
FROM (
  SELECT 'loan-module.products.view' AS pname UNION ALL
  SELECT 'loan-module.products.create' UNION ALL
  SELECT 'loan-module.products.update' UNION ALL
  SELECT 'loan-module.applications.view' UNION ALL
  SELECT 'loan-module.applications.create' UNION ALL
  SELECT 'loan-module.applications.submit.create' UNION ALL
  SELECT 'loan-module.applications.approve.create' UNION ALL
  SELECT 'loan-module.applications.reject.create' UNION ALL
  SELECT 'loan-module.contracts.view' UNION ALL
  SELECT 'loan-module.contracts.disburse.create' UNION ALL
  SELECT 'loan-module.contracts.repayments.create' UNION ALL
  SELECT 'loan-module.contracts.schedule.view' UNION ALL
  SELECT 'loan-module.contracts.statement.view' UNION ALL
  SELECT 'loan-module.contracts.settle.create' UNION ALL
  SELECT 'loan-module.contracts.restructure.create' UNION ALL
  SELECT 'loan-module.contracts.write-off.create' UNION ALL
  SELECT 'loan-module.disbursements.confirm.create' UNION ALL
  SELECT 'loan-module.disbursements.provider-status.view' UNION ALL
  SELECT 'loan-module.interest.accrue.create' UNION ALL
  SELECT 'loan-module.penalties.run.create' UNION ALL
  SELECT 'loan-module.reports.view'
) AS seed;

UPDATE roles SET deleted_at = NULL WHERE name = 'Super Admin';

INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.name LIKE 'loan-module.%'
WHERE r.name = 'Super Admin';
