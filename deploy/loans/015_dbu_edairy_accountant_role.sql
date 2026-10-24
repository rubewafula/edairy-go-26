-- DBU Edairy Accountant: the only role allowed to clear (settle) a month's dairy loans.
-- The API also creates these on startup (EnsureSettlementPermission).
INSERT INTO permissions (name, guard_name, created_at, updated_at)
SELECT 'loan-module.dairy-loans.clear', 'web', NOW(), NOW() FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE name = 'loan-module.dairy-loans.clear');

INSERT INTO roles (name, guard_name, created_at, updated_at)
SELECT 'DBU Edairy Accountant', 'web', NOW(), NOW() FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM roles WHERE name = 'DBU Edairy Accountant');

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r
JOIN permissions p ON p.deleted_at IS NULL AND p.name IN (
  'loan-module.dairy-loans.clear',
  'loan-module.applications.settle-month.create',
  'loan-module.applications.settlement-summary.view',
  'loan-module.applications.view',
  'loan-module.products.view')
WHERE r.name = 'DBU Edairy Accountant'
AND NOT EXISTS (SELECT 1 FROM role_permissions x WHERE x.role_id = r.id AND x.permission_id = p.id);
