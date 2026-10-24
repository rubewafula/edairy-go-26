-- Permission for GET /api/milk-collection-overview (mobile dashboard "Milk Delivery Overview").
-- RequireAutoPermission maps that route to "milk-collection-overview.view".
-- Grants it to every role that can already view milk journals. Safe to re-run.

INSERT IGNORE INTO permissions (name, guard_name, created_at, updated_at)
VALUES ('milk-collection-overview.view', 'web', NOW(), NOW());

INSERT IGNORE INTO role_permissions (permission_id, role_id, created_at, updated_at)
SELECT newp.id, rp.role_id, NOW(), NOW()
FROM permissions newp
JOIN permissions src ON src.name = 'milk-journals.view' AND src.guard_name = 'web' AND src.deleted_at IS NULL
JOIN role_permissions rp ON rp.permission_id = src.id AND rp.deleted_at IS NULL
WHERE newp.name = 'milk-collection-overview.view' AND newp.guard_name = 'web';
