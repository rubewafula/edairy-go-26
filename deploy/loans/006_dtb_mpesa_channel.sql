-- DTB M-Pesa disbursement channel + confirm/disburse permissions
INSERT INTO disbursement_channels
  (channel_code, channel_name, description, channel_type, provider, is_async, requires_external_ref, gl_rule_type, sort_order, is_active, config_schema, created_at, updated_at)
VALUES
  ('MOBILE_DTB_MPESA', 'Mobile - DTB M-Pesa', 'DTB/Astra M-Pesa payout (KE_DTB_MPESA) with OTP confirmation', 'MOBILE', 'DTB', 1, 1, 'LOAN_DISBURSEMENT', 35, 1,
   JSON_OBJECT('required', JSON_ARRAY('phone_number')), NOW(), NOW())
ON DUPLICATE KEY UPDATE
  channel_name = VALUES(channel_name),
  description = VALUES(description),
  provider = VALUES(provider),
  is_async = VALUES(is_async),
  requires_external_ref = VALUES(requires_external_ref),
  config_schema = VALUES(config_schema),
  is_active = VALUES(is_active),
  updated_at = NOW();

INSERT IGNORE INTO permissions (name, guard_name, created_at, updated_at)
SELECT pname, 'web', NOW(), NOW()
FROM (
  SELECT 'loan-module.contracts.disburse.create' AS pname UNION ALL
  SELECT 'loan-module.disbursements.confirm.create'
) AS seed;

INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.name IN (
  'loan-module.contracts.disburse.create',
  'loan-module.disbursements.confirm.create'
)
WHERE r.name = 'Super Admin';

INSERT IGNORE INTO user_permissions (user_id, permission_id)
SELECT u.id, p.id
FROM users u
JOIN permissions p ON p.name IN (
  'loan-module.contracts.disburse.create',
  'loan-module.disbursements.confirm.create'
)
WHERE u.email = 'rubewafula@gmail.com';
