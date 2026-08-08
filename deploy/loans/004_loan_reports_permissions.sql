-- Loan module report endpoint permissions (RBAC maps path segments to permission names)
INSERT IGNORE INTO permissions (name, guard_name, created_at, updated_at)
SELECT pname, 'web', NOW(), NOW()
FROM (
  SELECT 'loan-module.reports.view' AS pname UNION ALL
  SELECT 'loan-module.reports.portfolio.view' UNION ALL
  SELECT 'loan-module.reports.aging.view' UNION ALL
  SELECT 'loan-module.reports.par.view' UNION ALL
  SELECT 'loan-module.reports.register.view' UNION ALL
  SELECT 'loan-module.reports.applications.view' UNION ALL
  SELECT 'loan-module.reports.disbursements.view'
) AS seed;

INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.name IN (
  'loan-module.reports.view',
  'loan-module.reports.portfolio.view',
  'loan-module.reports.aging.view',
  'loan-module.reports.par.view',
  'loan-module.reports.register.view',
  'loan-module.reports.applications.view',
  'loan-module.reports.disbursements.view'
)
WHERE r.name = 'Super Admin';

INSERT IGNORE INTO user_permissions (user_id, permission_id)
SELECT u.id, p.id
FROM users u
JOIN permissions p ON p.name IN (
  'loan-module.reports.view',
  'loan-module.reports.portfolio.view',
  'loan-module.reports.aging.view',
  'loan-module.reports.par.view',
  'loan-module.reports.register.view',
  'loan-module.reports.applications.view',
  'loan-module.reports.disbursements.view'
)
WHERE u.email = 'rubewafula@gmail.com';
