-- Disbursement channel provider config permissions
INSERT IGNORE INTO permissions (name, guard_name, created_at, updated_at)
SELECT pname, 'web', NOW(), NOW()
FROM (
  SELECT 'loan-module.disbursement-channels.config.view' AS pname UNION ALL
  SELECT 'loan-module.disbursement-channels.config.update'
) AS seed;

INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.name IN (
  'loan-module.disbursement-channels.config.view',
  'loan-module.disbursement-channels.config.update'
)
WHERE r.name = 'Super Admin';

INSERT IGNORE INTO user_permissions (user_id, permission_id)
SELECT u.id, p.id
FROM users u
JOIN permissions p ON p.name IN (
  'loan-module.disbursement-channels.config.view',
  'loan-module.disbursement-channels.config.update'
)
WHERE u.email = 'rubewafula@gmail.com';
