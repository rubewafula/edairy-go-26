-- Member credit limit and mobile loan RBAC permissions (idempotent).
-- The API also seeds these at startup; this script is for databases updated by hand.

INSERT INTO permissions (name, guard_name, created_at, updated_at)
SELECT seed.pname, 'web', NOW(), NOW()
FROM (
  SELECT 'member-credit-limit.view' AS pname UNION ALL
  SELECT 'members.credit-limit.view' UNION ALL
  SELECT 'members.credit-limit-logs.view' UNION ALL
  SELECT 'members.credit-limit.recalculate.create' UNION ALL
  SELECT 'member-credit-limits.recalculate-all.create' UNION ALL
  SELECT 'loan-module.mobile.loan-requests.create' UNION ALL
  SELECT 'loan-module.mobile.loan-requests.view' UNION ALL
  SELECT 'loan-module.mobile.loan-requests.resend-otp.create' UNION ALL
  SELECT 'loan-module.mobile.loan-requests.verify-otp.create'
) AS seed
WHERE NOT EXISTS (SELECT 1 FROM permissions p WHERE p.name = seed.pname);

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.name IN (
  'member-credit-limit.view',
  'members.credit-limit.view',
  'members.credit-limit-logs.view',
  'members.credit-limit.recalculate.create',
  'member-credit-limits.recalculate-all.create',
  'loan-module.mobile.loan-requests.create',
  'loan-module.mobile.loan-requests.view',
  'loan-module.mobile.loan-requests.resend-otp.create',
  'loan-module.mobile.loan-requests.verify-otp.create'
)
WHERE r.name = 'Super Admin'
AND NOT EXISTS (SELECT 1 FROM role_permissions rp WHERE rp.role_id = r.id AND rp.permission_id = p.id);
