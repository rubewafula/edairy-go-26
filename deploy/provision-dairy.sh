#!/usr/bin/env bash
# Provision a new eDairy cooperative on edairy-db (systemd + nginx + HAProxy layout).
# Usage: sudo MYSQL_ROOT='...' ./provision-dairy.sh <slug> <api_port> <ui_port> [source_db]
set -euo pipefail

SLUG="${1:?slug required, e.g. mukululu}"
API_PORT="${2:?api port required}"
UI_PORT="${3:?ui port required}"
SOURCE_DB="${4:-arithi_dairy_erp}"
DB_NAME="${SLUG//-/_}_dairy_erp"
DB_USER="edairy"
APPS_ROOT="/apps/go"
NGINX_ROOT="/usr/share/nginx"
SERVICE_NAME="${SLUG}-erp"
BINARY_NAME="${SLUG}-erp-api"
DOMAIN_UI="${SLUG}.edairy.africa"
DOMAIN_API="api.${SLUG}.edairy.africa"
SLUG_VAR="${SLUG//-/_}"

if [[ "${EUID}" -ne 0 ]]; then
  echo "Run as root: sudo MYSQL_ROOT='...' $0 $*" >&2
  exit 1
fi

DB_PASS="$(grep DB_PASSWORD "${APPS_ROOT}/arithi/.env" | cut -d= -f2- | tr -d ' ')"
JWT_SECRET="$(openssl rand -hex 32)"

mysql_root() {
  if [[ -n "${MYSQL_ROOT:-}" ]]; then
    mysql -uroot -p"${MYSQL_ROOT}" "$@"
  else
    mysql "$@"
  fi
}

echo "==> Creating database ${DB_NAME} (schema only from ${SOURCE_DB})"
mysql_root <<SQL
CREATE DATABASE IF NOT EXISTS \`${DB_NAME}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON \`${DB_NAME}\`.* TO '${DB_USER}'@'localhost';
GRANT ALL PRIVILEGES ON \`${DB_NAME}\`.* TO '${DB_USER}'@'127.0.0.1';
FLUSH PRIVILEGES;
SQL

mysqldump $( [[ -n "${MYSQL_ROOT:-}" ]] && echo -uroot -p"${MYSQL_ROOT}" ) --no-data --routines --triggers "${SOURCE_DB}" \
  | mysql_root "${DB_NAME}"

echo "==> Backend directory ${APPS_ROOT}/${SLUG}"
install -d -m 775 "${APPS_ROOT}/${SLUG}/uploads/members" "${APPS_ROOT}/${SLUG}/uploads/transporters" "${APPS_ROOT}/${SLUG}/storage/exports"

cat > "${APPS_ROOT}/${SLUG}/.env" <<ENV
APP_NAME=${SLUG^^}
APP_ENV=production
PORT=${API_PORT}
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=${DB_NAME}
DB_USER=${DB_USER}
DB_PASSWORD=${DB_PASS}
ENV=production
JWT_SECRET=${JWT_SECRET}
JWT_EXPIRY_HOURS=24
ENV

chmod 640 "${APPS_ROOT}/${SLUG}/.env"
chown reuben:reuben "${APPS_ROOT}/${SLUG}/.env" "${APPS_ROOT}/${SLUG}" || true

if [[ -f "${APPS_ROOT}/${SLUG}/${BINARY_NAME}" ]]; then
  chmod +x "${APPS_ROOT}/${SLUG}/${BINARY_NAME}"
fi

echo "==> Nginx frontend ${NGINX_ROOT}/${SLUG} (port ${UI_PORT})"
install -d -m 755 "${NGINX_ROOT}/${SLUG}"
chown reuben:reuben "${NGINX_ROOT}/${SLUG}" || true

cat > "/etc/nginx/conf.d/${SLUG}.conf" <<NGINX
server {
    listen ${UI_PORT};
    server_name ${DOMAIN_UI};

    root ${NGINX_ROOT}/${SLUG};
    index index.html;

    location / {
        try_files \$uri /index.html;
    }

    location /static/ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    location ~* \\.(js|css|png|jpg|jpeg|gif|svg|ico|woff|woff2|ttf)\$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
        access_log off;
    }

    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/javascript application/javascript application/json application/xml image/svg+xml;

    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    add_header Permissions-Policy "geolocation=(), microphone=(), camera=()" always;

    access_log /var/log/nginx/${SLUG}-access.log;
    error_log /var/log/nginx/${SLUG}-error.log;
}
NGINX

echo "==> systemd unit ${SERVICE_NAME}.service"
cat > "/etc/systemd/system/${SERVICE_NAME}.service" <<UNIT
[Unit]
Description=${SLUG} ERP Go API Service
After=network.target mysql.service
Wants=mysql.service

[Service]
User=root
WorkingDirectory=${APPS_ROOT}/${SLUG}
ExecStart=${APPS_ROOT}/${SLUG}/${BINARY_NAME}
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
UNIT

systemctl daemon-reload
systemctl enable "${SERVICE_NAME}.service"

HAPROXY_CFG="/etc/haproxy/haproxy.cfg"
if ! grep -q "host_${SLUG_VAR}_ui_live" "$HAPROXY_CFG"; then
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
fi

nginx -t
systemctl reload nginx

echo "Done: ${SLUG} (db=${DB_NAME}, api=${API_PORT}, ui=${UI_PORT})"
echo "Deploy binary to ${APPS_ROOT}/${SLUG}/${BINARY_NAME} and frontend to ${NGINX_ROOT}/${SLUG}/"
echo "Then: systemctl start ${SERVICE_NAME}.service"
