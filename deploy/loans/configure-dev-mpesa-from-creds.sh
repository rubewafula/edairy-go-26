#!/bin/bash
# Configure dev stack M-Pesa from a local creds folder, then recreate backend.
# Expected layout under CREDS_DIR (default: /home/reuben/apps/mpesa.creds):
#   mpesa.env          — MPESA_* variables (see edairy-go/.env.example)
#   production.cer     — Safaricom public cert (or any *.cer; first match used)
#
# Callback URLs are forced to api.dev.edairy.africa (Safaricom must reach these).
set -euo pipefail

CREDS_PATH="${1:-/home/reuben/apps/mpesa.creds}"
CERT_PATH="${2:-${MPESA_CERT:-/home/reuben/apps/mpesa.cer}}"
DEV_STACK="${DEV_STACK:-/opt/edairy/dev}"
REMOTE="${REMOTE:-edairy-db}"

ENV_FILE=""
if [[ -f "$CREDS_PATH" ]]; then
  ENV_FILE="$CREDS_PATH"
elif [[ -d "$CREDS_PATH" ]]; then
  ENV_FILE="$CREDS_PATH/mpesa.env"
  if [[ ! -f "$ENV_FILE" ]]; then
    for f in "$CREDS_PATH"/.env "$CREDS_PATH"/env; do
      [[ -f "$f" ]] && ENV_FILE="$f" && break
    done
  fi
else
  echo "Missing creds file or directory: $CREDS_PATH" >&2
  exit 1
fi
if [[ ! -f "$ENV_FILE" ]]; then
  echo "Missing mpesa.env under $CREDS_PATH" >&2
  exit 1
fi

if [[ ! -f "$CERT_PATH" ]]; then
  echo "Missing certificate: $CERT_PATH" >&2
  exit 1
fi
CERT="$CERT_PATH"

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT
# Normalize initiator password aliases (password= / initiator_password=)
python3 - "$ENV_FILE" "$TMP/mpesa.env" <<'PY'
import re, sys, pathlib
src, dst = pathlib.Path(sys.argv[1]), pathlib.Path(sys.argv[2])
lines = src.read_text().splitlines()
out, pwd = [], None
label_map = {
    'user name': 'MPESA_INITIATOR_NAME',
    'pay bill number': 'MPESA_B2C_SHORTCODE',
    'password': 'MPESA_INITIATOR_PASSWORD',
    'consumer key': 'MPESA_CONSUMER_KEY',
    'consumer secret': 'MPESA_CONSUMER_SECRET',
}
for line in lines:
    s = line.strip()
    if not s or s.startswith('#'):
        continue
    if s.startswith('password=') or s.startswith('initiator_password='):
        pwd = s.split('=', 1)[1].strip().strip('"').strip("'")
        continue
    if ':' in s and not s.startswith('MPESA_'):
        label, val = s.split(':', 1)
        norm = re.sub(r'^\d+\.\s*', '', label.strip()).lower()
        if norm in label_map:
            out.append(f'{label_map[norm]}={val.strip()}')
            continue
    if '=' in s and not s.startswith('MPESA_'):
        k, v = s.split('=', 1)
        k = k.strip().upper()
        mapping = {
            'CONSUMER_KEY': 'MPESA_CONSUMER_KEY',
            'CONSUMER_SECRET': 'MPESA_CONSUMER_SECRET',
            'B2C_SHORTCODE': 'MPESA_B2C_SHORTCODE',
            'SHORTCODE': 'MPESA_B2C_SHORTCODE',
            'INITIATOR_NAME': 'MPESA_INITIATOR_NAME',
            'INITIATOR_PASSWORD': 'MPESA_INITIATOR_PASSWORD',
            'ENV': 'MPESA_ENV',
        }
        if k in mapping:
            out.append(f'{mapping[k]}={v.strip()}')
            continue
    out.append(s)
if pwd and not any(x.startswith('MPESA_INITIATOR_PASSWORD=') for x in out):
    out.append(f'MPESA_INITIATOR_PASSWORD={pwd}')
if not any(x.startswith('MPESA_ENV=') for x in out):
    out.insert(0, 'MPESA_ENV=production')
dst.write_text('\n'.join(out) + '\n')
PY
cp "$CERT" "$TMP/mpesa-cert.cer"

ssh "$REMOTE" "echo '0tuma@01' | sudo -S mkdir -p $DEV_STACK/data/mpesa && echo '0tuma@01' | sudo -S chown reuben:reuben $DEV_STACK/data/mpesa"
scp "$TMP/mpesa-cert.cer" "$REMOTE:/tmp/mpesa-cert.cer"
scp "$TMP/mpesa.env" "$REMOTE:/tmp/mpesa.env.partial"

ssh "$REMOTE" bash -s <<'EOS'
set -euo pipefail
SUDO_PW='0tuma@01'
sudo() { echo "$SUDO_PW" | command sudo -S "$@"; }
DEV_STACK="/opt/edairy/dev"
sudo mv /tmp/mpesa-cert.cer "$DEV_STACK/data/mpesa/production.cer"
sudo chown root:root "$DEV_STACK/data/mpesa/production.cer"
sudo chmod 644 "$DEV_STACK/data/mpesa/production.cer"

# Merge MPESA_* into existing .env (never truncate base DB/JWT settings)
if ! sudo test -s "$DEV_STACK/.env"; then
  echo "ERROR: $DEV_STACK/.env is missing or empty; restore DB settings before M-Pesa configure." >&2
  exit 1
fi
if ! sudo test -s /tmp/mpesa.env.partial; then
  echo "ERROR: /tmp/mpesa.env.partial missing or empty (scp failed?)" >&2
  exit 1
fi
sudo grep -v '^MPESA_' "$DEV_STACK/.env" | sudo tee /tmp/dev.env.base >/dev/null
if ! sudo grep -q '^DB_NAME=' /tmp/dev.env.base; then
  echo "ERROR: base .env missing DB settings; refusing to overwrite $DEV_STACK/.env" >&2
  exit 1
fi
{
  sudo cat /tmp/dev.env.base
  sudo cat /tmp/mpesa.env.partial
  echo "MPESA_CERT_PATH=/app/mpesa/production.cer"
  echo "MPESA_B2C_RESULT_URL=https://api.dev.edairy.africa/api/mpesa-b2c-result"
  echo "MPESA_B2C_TIMEOUT_URL=https://api.dev.edairy.africa/api/mpesa-b2c-timeout"
} | sudo tee "$DEV_STACK/.env" >/dev/null
sudo chmod 640 "$DEV_STACK/.env"
sudo rm -f /tmp/mpesa.env.partial /tmp/dev.env.base
if ! sudo test -s "$DEV_STACK/.env"; then
  echo "ERROR: merged .env is empty" >&2
  exit 1
fi
if ! sudo grep -q '^DB_NAME=' "$DEV_STACK/.env"; then
  echo "ERROR: merged .env lost DB settings" >&2
  exit 1
fi

COMPOSE="$DEV_STACK/docker-compose.yml"
if ! sudo grep -q '/app/mpesa' "$COMPOSE"; then
  sudo sed -i '/- \.\/data\/storage:\/app\/storage/a\      - ./data/mpesa:/app/mpesa:ro' "$COMPOSE"
fi

cd "$DEV_STACK"
sudo docker compose up -d --force-recreate backend
sudo docker compose ps backend
EOS

echo "Dev backend recreated with M-Pesa env. Run dev-mpesa-loan-disburse-test.sh next."
