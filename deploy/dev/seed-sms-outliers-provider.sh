#!/bin/bash
# Upsert default Outliers SMS provider on dev from SMS_OUTLIERS_* in /opt/edairy/dev/.env
set -euo pipefail
ENV_FILE="${ENV_FILE:-/opt/edairy/dev/.env}"
[[ -f "$ENV_FILE" ]] || { echo "Missing $ENV_FILE" >&2; exit 1; }
env_get() { grep -m1 "^$1=" "$ENV_FILE" | cut -d= -f2-; }
SMS_OUTLIERS_CLIENT_CODE="$(env_get SMS_OUTLIERS_CLIENT_CODE)"
SMS_OUTLIERS_API_KEY="$(env_get SMS_OUTLIERS_API_KEY)"
SMS_OUTLIERS_SENDER_ID="$(env_get SMS_OUTLIERS_SENDER_ID)"
SMS_OUTLIERS_API_URL="$(env_get SMS_OUTLIERS_API_URL)"
DB_PASSWORD="$(env_get DB_PASSWORD)"
DB_USER="$(env_get DB_USER)"
DB_NAME="$(env_get DB_NAME)"
: "${SMS_OUTLIERS_CLIENT_CODE:?Set SMS_OUTLIERS_CLIENT_CODE in .env}"
: "${SMS_OUTLIERS_API_KEY:?Set SMS_OUTLIERS_API_KEY in .env}"
SMS_OUTLIERS_SENDER_ID="${SMS_OUTLIERS_SENDER_ID:-E-Dairy}"
SMS_OUTLIERS_API_URL="${SMS_OUTLIERS_API_URL:-https://sms.outliers.co.ke/api/v2/send-sms}"
DEV_PASS="${DB_PASSWORD:?DB_PASSWORD missing in .env}"
mysql -h 127.0.0.1 -u"${DB_USER:-edairy}" -p"$DEV_PASS" "${DB_NAME:-dev_dairy_erp}" <<SQL
UPDATE sms_providers SET is_default='no' WHERE provider_code != '${SMS_OUTLIERS_CLIENT_CODE}';
INSERT INTO sms_providers (
  provider_code, provider_name, api_url, api_key, sender_id, is_default, status, created_at, updated_at, created_by, updated_by
) VALUES (
  '${SMS_OUTLIERS_CLIENT_CODE}', 'Outliers', '${SMS_OUTLIERS_API_URL}', '${SMS_OUTLIERS_API_KEY}',
  '${SMS_OUTLIERS_SENDER_ID}', 'yes', 'active', NOW(), NOW(), 1, 1
)
ON DUPLICATE KEY UPDATE
  provider_name='Outliers',
  api_url=VALUES(api_url),
  api_key=VALUES(api_key),
  sender_id=VALUES(sender_id),
  is_default='yes',
  status='active',
  updated_at=NOW(),
  updated_by=1;
SQL
echo "Default Outliers SMS provider ${SMS_OUTLIERS_CLIENT_CODE} (${SMS_OUTLIERS_SENDER_ID}) ready."
