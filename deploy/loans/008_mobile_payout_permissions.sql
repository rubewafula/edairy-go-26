-- Mobile payout disburse + confirm + provider status permissions
INSERT IGNORE INTO permissions (name, guard_name, created_at, updated_at)
SELECT pname, 'web', NOW(), NOW()
FROM (
  SELECT 'loan-module.contracts.disburse.create' AS pname UNION ALL
  SELECT 'loan-module.disbursements.confirm.create' UNION ALL
  SELECT 'loan-module.disbursements.provider-status.view'
) AS seed;

INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.name IN (
  'loan-module.contracts.disburse.create',
  'loan-module.disbursements.confirm.create',
  'loan-module.disbursements.provider-status.view'
)
WHERE r.name = 'Super Admin';

INSERT IGNORE INTO user_permissions (user_id, permission_id)
SELECT u.id, p.id
FROM users u
JOIN permissions p ON p.name IN (
  'loan-module.contracts.disburse.create',
  'loan-module.disbursements.confirm.create',
  'loan-module.disbursements.provider-status.view'
)
WHERE u.email = 'rubewafula@gmail.com';
