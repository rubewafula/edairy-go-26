#!/usr/bin/env bash
# Rename a dairy stack on edairy-db (Docker, MySQL, HAProxy).
# Usage: sudo MYSQL_ROOT='...' ./rename-dairy.sh <old_slug> <new_slug> [api_port] [ui_port]
#
# Example: sudo ./rename-dairy.sh nguene nkuene 8510 8511
set -euo pipefail

OLD_SLUG="${1:?old slug required}"
NEW_SLUG="${2:?new slug required}"
API_PORT="${3:-}"
UI_PORT="${4:-}"
OLD_VAR="${OLD_SLUG//-/_}"
NEW_VAR="${NEW_SLUG//-/_}"
OLD_DB="${OLD_SLUG//-/_}_dairy_erp"
NEW_DB="${NEW_SLUG//-/_}_dairy_erp"
OLD_STACK="/opt/edairy/${OLD_SLUG}"
NEW_STACK="/opt/edairy/${NEW_SLUG}"
HAPROXY_CFG="/etc/haproxy/haproxy.cfg"

if [[ "${EUID}" -ne 0 ]]; then
  echo "Run as root: sudo $0 $*" >&2
  exit 1
fi

if [[ ! -d "${OLD_STACK}" ]]; then
  echo "Stack not found: ${OLD_STACK}" >&2
  exit 1
fi

if [[ -z "${API_PORT}" ]]; then
  API_PORT="$(grep '^PORT=' "${OLD_STACK}/.env" | cut -d= -f2-)"
fi
if [[ -z "${UI_PORT}" ]]; then
  UI_PORT="$(grep -E '^\s+- "[0-9]+:80"' "${OLD_STACK}/docker-compose.yml" | head -1 | sed 's/.*"\([0-9]*\):80".*/\1/')"
fi

mysql_root() {
  if [[ -n "${MYSQL_ROOT:-}" ]]; then
    mysql -uroot -p"${MYSQL_ROOT}" "$@"
  else
    mysql "$@"
  fi
}

echo "==> Stopping ${OLD_SLUG} stack"
cd "${OLD_STACK}"
docker compose down

echo "==> Renaming database ${OLD_DB} -> ${NEW_DB}"
mysql_root <<SQL
CREATE DATABASE IF NOT EXISTS \`${NEW_DB}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON \`${NEW_DB}\`.* TO 'edairy'@'localhost';
GRANT ALL PRIVILEGES ON \`${NEW_DB}\`.* TO 'edairy'@'127.0.0.1';
GRANT ALL PRIVILEGES ON \`${NEW_DB}\`.* TO 'edairy'@'172.%';
FLUSH PRIVILEGES;
SQL

mysqldump $( [[ -n "${MYSQL_ROOT:-}" ]] && echo -uroot -p"${MYSQL_ROOT}" ) \
  --routines --triggers "${OLD_DB}" \
  | mysql_root "${NEW_DB}"

echo "==> Moving stack ${OLD_STACK} -> ${NEW_STACK}"
mv "${OLD_STACK}" "${NEW_STACK}"

echo "==> Updating stack configuration"
sed -i \
  -e "s/${OLD_SLUG}/${NEW_SLUG}/g" \
  -e "s/${OLD_VAR}/${NEW_VAR}/g" \
  -e "s/${OLD_DB}/${NEW_DB}/g" \
  -e "s/Nguene/Nkuene/g" \
  -e "s/NGUENE/NKUENE/g" \
  "${NEW_STACK}/.env" "${NEW_STACK}/docker-compose.yml" "${NEW_STACK}/runtime-config.js"

echo "==> Updating HAProxy (${OLD_SLUG} -> ${NEW_SLUG})"
python3 - <<PY
from pathlib import Path
path = Path("${HAPROXY_CFG}")
text = path.read_text()
replacements = [
    ("host_${OLD_VAR}_ui_live", "host_${NEW_VAR}_ui_live"),
    ("host_${OLD_VAR}_api_live", "host_${NEW_VAR}_api_live"),
    ("${OLD_VAR}_erp_backend", "${NEW_VAR}_erp_backend"),
    ("${OLD_VAR}_erp_frontend_backend", "${NEW_VAR}_erp_frontend_backend"),
    ("${OLD_SLUG}.edairy.africa", "${NEW_SLUG}.edairy.africa"),
    ("api.${OLD_SLUG}.edairy.africa", "api.${NEW_SLUG}.edairy.africa"),
    ("${OLD_VAR}_ui_live", "${NEW_VAR}_ui_live"),
    ("${OLD_VAR}_live_erp", "${NEW_VAR}_live_erp"),
    ("backend ${OLD_VAR}_erp_frontend_backend", "backend ${NEW_VAR}_erp_frontend_backend"),
    ("backend ${OLD_VAR}_erp_backend", "backend ${NEW_VAR}_erp_backend"),
]
for old, new in replacements:
    text = text.replace(old, new)
path.write_text(text)
PY

haproxy -c -f "${HAPROXY_CFG}"
systemctl reload haproxy

echo "==> Starting ${NEW_SLUG} stack"
cd "${NEW_STACK}"
docker compose up -d --pull never

echo "==> Dropping old database ${OLD_DB}"
mysql_root -e "DROP DATABASE IF EXISTS \`${OLD_DB}\`;"

echo "Done: ${OLD_SLUG} -> ${NEW_SLUG} (db=${NEW_DB}, api=${API_PORT}, ui=${UI_PORT})"
