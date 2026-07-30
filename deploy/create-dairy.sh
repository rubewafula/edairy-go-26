#!/usr/bin/env bash
# Provision or migrate an eDairy cooperative using Docker Compose on edairy-db.
#
# Usage:
#   sudo MYSQL_ROOT='...' ./create-dairy.sh [--migrate] [--skip-db] [--skip-haproxy] [--no-start] \
#       <slug> <api_port> <ui_port> [source_db]
#
# Examples:
#   # New dairy
#   sudo MYSQL_ROOT='...' ./create-dairy.sh mukululu 8504 8505
#
#   # Migrate existing systemd deployment to Docker
#   sudo ./create-dairy.sh --migrate --skip-db arithi 8500 8501
#
#   # Migrate all six dairies
#   sudo ./migrate-all-to-docker.sh
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATES="${SCRIPT_DIR}/templates"
EDAIRY_ROOT="/opt/edairy"
LEGACY_APPS="/apps/go"
EDAIRY_VERSION="${EDAIRY_VERSION:-1.0.0}"

MIGRATE=false
SKIP_DB=false
SKIP_HAPROXY=false
NO_START=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --migrate) MIGRATE=true; shift ;;
    --skip-db) SKIP_DB=true; shift ;;
    --skip-haproxy) SKIP_HAPROXY=true; shift ;;
    --no-start) NO_START=true; shift ;;
    -*) echo "Unknown option: $1" >&2; exit 1 ;;
    *) break ;;
  esac
done

SLUG="${1:?slug required, e.g. mukululu}"
API_PORT="${2:?api port required}"
UI_PORT="${3:?ui port required}"
SOURCE_DB="${4:-arithi_dairy_erp}"
DB_NAME="${SLUG//-/_}_dairy_erp"
DB_USER="edairy"
DOMAIN_UI="${SLUG}.edairy.africa"
DOMAIN_API="api.${SLUG}.edairy.africa"
SLUG_VAR="${SLUG//-/_}"
STACK_DIR="${EDAIRY_ROOT}/${SLUG}"
SERVICE_NAME="${SLUG}-erp"
APP_NAME="$(echo "${SLUG}" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) substr($i,2)}1') ERP"

if [[ "${EUID}" -ne 0 ]]; then
  echo "Run as root: sudo MYSQL_ROOT='...' $0 $*" >&2
  exit 1
fi

mysql_root() {
  if [[ -n "${MYSQL_ROOT:-}" ]]; then
    mysql -uroot -p"${MYSQL_ROOT}" "$@"
  else
    mysql "$@"
  fi
}

render_template() {
  local src="$1" dest="$2"
  sed \
    -e "s|{{SLUG}}|${SLUG}|g" \
    -e "s|{{APP_NAME}}|${APP_NAME}|g" \
    -e "s|{{API_PORT}}|${API_PORT}|g" \
    -e "s|{{UI_PORT}}|${UI_PORT}|g" \
    -e "s|{{DB_NAME}}|${DB_NAME}|g" \
    -e "s|{{DB_PASSWORD}}|${DB_PASS}|g" \
    -e "s|{{JWT_SECRET}}|${JWT_SECRET}|g" \
    -e "s|{{EDAIRY_VERSION}}|${EDAIRY_VERSION}|g" \
    "$src" > "$dest"
}

resolve_secrets() {
  if [[ -f "${LEGACY_APPS}/${SLUG}/.env" ]]; then
    DB_PASS="$(grep '^DB_PASSWORD=' "${LEGACY_APPS}/${SLUG}/.env" | cut -d= -f2- | tr -d ' ')"
    JWT_SECRET="$(grep '^JWT_SECRET=' "${LEGACY_APPS}/${SLUG}/.env" | cut -d= -f2- | tr -d ' ')"
  elif [[ -f "${STACK_DIR}/.env" ]]; then
    DB_PASS="$(grep '^DB_PASSWORD=' "${STACK_DIR}/.env" | cut -d= -f2- | tr -d ' ')"
    JWT_SECRET="$(grep '^JWT_SECRET=' "${STACK_DIR}/.env" | cut -d= -f2- | tr -d ' ')"
  elif [[ -f "${LEGACY_APPS}/arithi/.env" ]]; then
    DB_PASS="$(grep '^DB_PASSWORD=' "${LEGACY_APPS}/arithi/.env" | cut -d= -f2- | tr -d ' ')"
    JWT_SECRET="$(openssl rand -hex 32)"
  else
    echo "Cannot resolve DB_PASSWORD; set LEGACY_APPS/arithi/.env or pass existing stack .env" >&2
    exit 1
  fi

  [[ -n "${DB_PASS}" ]] || { echo "DB_PASSWORD is empty" >&2; exit 1; }
  [[ -n "${JWT_SECRET}" ]] || JWT_SECRET="$(openssl rand -hex 32)"
}

install_templates() {
  install -d -m 755 "${EDAIRY_ROOT}/templates"
  install -m 644 "${TEMPLATES}/docker-compose.yml" "${EDAIRY_ROOT}/templates/"
  install -m 644 "${TEMPLATES}/.env.template" "${EDAIRY_ROOT}/templates/"
  install -m 644 "${TEMPLATES}/runtime-config.js.template" "${EDAIRY_ROOT}/templates/"
  install -d -m 755 "${EDAIRY_ROOT}/templates/nginx"
  install -m 644 "${TEMPLATES}/nginx/default.conf" "${EDAIRY_ROOT}/templates/nginx/"
}

provision_database() {
  echo "==> Creating database ${DB_NAME} (schema only from ${SOURCE_DB})"
  mysql_root <<SQL
CREATE DATABASE IF NOT EXISTS \`${DB_NAME}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON \`${DB_NAME}\`.* TO '${DB_USER}'@'localhost';
GRANT ALL PRIVILEGES ON \`${DB_NAME}\`.* TO '${DB_USER}'@'127.0.0.1';
GRANT ALL PRIVILEGES ON \`${DB_NAME}\`.* TO '${DB_USER}'@'172.%';
FLUSH PRIVILEGES;
SQL

  mysqldump $( [[ -n "${MYSQL_ROOT:-}" ]] && echo -uroot -p"${MYSQL_ROOT}" ) \
    --no-data --routines --triggers "${SOURCE_DB}" \
    | mysql_root "${DB_NAME}"
}

ensure_docker_grants() {
  echo "==> Ensuring MySQL grants for Docker bridge (172.*)"
  mysql_root <<SQL
GRANT ALL PRIVILEGES ON \`${DB_NAME}\`.* TO '${DB_USER}'@'172.%';
FLUSH PRIVILEGES;
SQL
}

copy_legacy_data() {
  local legacy="${LEGACY_APPS}/${SLUG}"
  if [[ ! -d "${legacy}" ]]; then
    echo "==> No legacy data at ${legacy}; creating empty upload dirs"
    install -d -m 775 "${STACK_DIR}/data/uploads/members" \
      "${STACK_DIR}/data/uploads/transporters" \
      "${STACK_DIR}/data/storage/exports"
    return
  fi

  echo "==> Copying legacy uploads and storage from ${legacy}"
  if [[ -d "${legacy}/uploads" ]]; then
    rsync -a "${legacy}/uploads/" "${STACK_DIR}/data/uploads/"
  else
    install -d -m 775 "${STACK_DIR}/data/uploads/members" "${STACK_DIR}/data/uploads/transporters"
  fi

  if [[ -d "${legacy}/storage" ]]; then
    rsync -a "${legacy}/storage/" "${STACK_DIR}/data/storage/"
  else
    install -d -m 775 "${STACK_DIR}/data/storage/exports"
  fi
}

render_stack() {
  echo "==> Rendering Docker stack at ${STACK_DIR}"
  install -d -m 755 "${STACK_DIR}/data/uploads" "${STACK_DIR}/data/storage"

  export SLUG API_PORT UI_PORT EDAIRY_VERSION
  render_template "${TEMPLATES}/.env.template" "${STACK_DIR}/.env"
  render_template "${TEMPLATES}/runtime-config.js.template" "${STACK_DIR}/runtime-config.js"
  render_template "${TEMPLATES}/docker-compose.yml" "${STACK_DIR}/docker-compose.yml"

  chmod 640 "${STACK_DIR}/.env"
  chmod 644 "${STACK_DIR}/runtime-config.js" "${STACK_DIR}/docker-compose.yml"
}

patch_haproxy() {
  local HAPROXY_CFG="/etc/haproxy/haproxy.cfg"
  if grep -q "host_${SLUG_VAR}_ui_live" "$HAPROXY_CFG"; then
    echo "==> HAProxy already configured for ${SLUG}"
    return
  fi

  echo "==> HAProxy ACL/backends for ${DOMAIN_UI} and ${DOMAIN_API}"
  python3 - <<PY
from pathlib import Path
slug_var = "${SLUG_VAR}"
api_port = "${API_PORT}"
ui_port = "${UI_PORT}"
domain_ui = "${DOMAIN_UI}"
domain_api = "${DOMAIN_API}"
path = Path("${HAPROXY_CFG}")
text = path.read_text()
needle = "    acl host_tigania_west_api_live hdr(host) -i api.tigania-west.edairy.africa"
insert = needle + f"\n    acl host_{slug_var}_ui_live hdr(host) -i {domain_ui}\n    acl host_{slug_var}_api_live hdr(host) -i {domain_api}"
if needle not in text:
    raise SystemExit("HAProxy template anchor not found")
text = text.replace(needle, insert, 1)
needle = "    use_backend tigania_west_erp_backend if host_tigania_west_api_live"
insert = needle + f"\n    use_backend {slug_var}_erp_backend if host_{slug_var}_api_live\n    use_backend {slug_var}_erp_frontend_backend if host_{slug_var}_ui_live"
text = text.replace(needle, insert, 1)
backends = f"""
backend {slug_var}_erp_frontend_backend
    option redispatch
    option prefer-last-server
    retries 5
    balance source
    hash-type consistent
    http-request redirect scheme https unless {{ ssl_fc }}
    http-request set-header X-Forwarded-Proto https
    server {slug_var}_ui_live 127.0.0.1:{ui_port} check maxconn 200

backend {slug_var}_erp_backend
    option redispatch
    option prefer-last-server
    retries 5
    balance source
    hash-type consistent
    option http-server-close
    option forwardfor
    http-reuse safe
    http-request redirect scheme https unless {{ ssl_fc }}
    http-request set-header X-Forwarded-Proto https
    server {slug_var}_live_erp 127.0.0.1:{api_port} check maxconn 200
"""
path.write_text(text.rstrip() + backends + "\n")
PY
  haproxy -c -f "$HAPROXY_CFG"
  systemctl reload haproxy
}

stop_legacy() {
  if systemctl is-active --quiet "${SERVICE_NAME}.service" 2>/dev/null; then
    echo "==> Stopping legacy systemd service ${SERVICE_NAME}"
    systemctl stop "${SERVICE_NAME}.service"
  fi
}

remove_legacy_nginx() {
  local nginx_conf="/etc/nginx/conf.d/${SLUG}.conf"
  if [[ -f "${nginx_conf}" ]]; then
    echo "==> Removing host nginx config ${nginx_conf} (freeing port ${UI_PORT})"
    mv "${nginx_conf}" "${nginx_conf}.bak.$(date +%Y%m%d%H%M%S)"
    nginx -t && systemctl reload nginx
  fi
}

disable_legacy() {
  if systemctl is-enabled --quiet "${SERVICE_NAME}.service" 2>/dev/null; then
    echo "==> Disabling legacy systemd service ${SERVICE_NAME}"
    systemctl disable "${SERVICE_NAME}.service"
  fi
}

start_stack() {
  echo "==> Starting Docker stack (EDAIRY_VERSION=${EDAIRY_VERSION})"
  cd "${STACK_DIR}"
  if [[ "${SKIP_PULL:-}" != "true" ]]; then
    docker compose pull
  fi
  docker compose up -d --pull never
}

health_check() {
  echo "==> Health checks"
  local retries=30
  for i in $(seq 1 "$retries"); do
    if curl -sf "http://127.0.0.1:${API_PORT}/health" >/dev/null; then
      echo "    API /health OK"
      break
    fi
    if [[ "$i" -eq "$retries" ]]; then
      echo "API health check failed after ${retries} attempts" >&2
      docker compose -f "${STACK_DIR}/docker-compose.yml" logs backend >&2 || true
      exit 1
    fi
    sleep 2
  done

  if curl -sf "http://127.0.0.1:${UI_PORT}/runtime-config.js" | grep -q "__EDAIRY_CONFIG__"; then
    echo "    UI runtime-config.js OK"
  else
    for i in $(seq 1 10); do
      sleep 2
      if curl -sf "http://127.0.0.1:${UI_PORT}/runtime-config.js" | grep -q "__EDAIRY_CONFIG__"; then
        echo "    UI runtime-config.js OK"
        return
      fi
    done
    echo "UI runtime-config.js check failed" >&2
    exit 1
  fi
}

# --- main ---

install_templates
resolve_secrets

if [[ "${SKIP_DB}" == false ]]; then
  provision_database
else
  ensure_docker_grants
fi

render_stack
copy_legacy_data

if [[ "${SKIP_HAPROXY}" == false ]]; then
  patch_haproxy
fi

if [[ "${MIGRATE}" == true ]]; then
  stop_legacy
  remove_legacy_nginx
fi

if [[ "${NO_START}" == false ]]; then
  start_stack
  health_check
fi

if [[ "${MIGRATE}" == true && "${NO_START}" == false ]]; then
  disable_legacy
fi

echo "Done: ${SLUG} (db=${DB_NAME}, api=${API_PORT}, ui=${UI_PORT}, stack=${STACK_DIR})"
