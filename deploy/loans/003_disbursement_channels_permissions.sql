-- Disbursement channels RBAC permissions
INSERT IGNORE INTO permissions (name, guard_name, created_at, updated_at)
SELECT pname, 'web', NOW(), NOW()
FROM (
  SELECT 'loan-module.disbursement-channels.view' AS pname
) AS seed;

-- Grant to Super Admin role (if present)
INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.name = 'loan-module.disbursement-channels.view'
WHERE r.name = 'Super Admin';
