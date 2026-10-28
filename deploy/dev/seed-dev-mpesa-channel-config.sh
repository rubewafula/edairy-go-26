#!/bin/bash
# Seed MOBILE_MPESA disbursement_channel_config on dev from /apps/go/mpesa.creds + dev .env URLs.
# Run on edairy-db (or via: ssh edairy-db 'bash -s' < seed-dev-mpesa-channel-config.sh)
set -euo pipefail

CREDS="${CREDS:-/apps/go/mpesa.creds}"
ENV_FILE="${ENV_FILE:-/opt/edairy/dev/.env}"
CHANNEL_ID="${CHANNEL_ID:-3}"

python3 - "$CREDS" "$ENV_FILE" "$CHANNEL_ID" <<'PY'
import re, sys, subprocess
from pathlib import Path

def parse_checklist(path):
    labels = {
        'user name': ('initiator_name', False),
        'pay bill number': ('b2c_shortcode', False),
        'password': ('initiator_password', True),
        'consumer key': ('consumer_key', True),
        'consumer secret': ('consumer_secret', True),
    }
    out = {'env': 'production'}
    for line in Path(path).read_text().splitlines():
        if ':' not in line:
            continue
        label, val = line.split(':', 1)
        norm = re.sub(r'^\d+\.\s*', '', label.strip()).lower()
        if norm in labels:
            key, secret = labels[norm]
            out[key] = val.strip()
    return out

def load_dotenv(path):
    d = {}
    for line in Path(path).read_text().splitlines():
        s = line.strip()
        if not s or s.startswith('#') or '=' not in s:
            continue
        k, v = s.split('=', 1)
        d[k.strip()] = v.strip()
    return d

creds_path, env_path, channel_id = sys.argv[1:4]
vals = parse_checklist(creds_path)
dot = load_dotenv(env_path)
for k in ('b2c_result_url', 'b2c_timeout_url', 'cert_path'):
    env_key = 'MPESA_' + k.upper()
    if dot.get(env_key):
        vals[k] = dot[env_key]
vals.setdefault('b2c_command_id', 'BusinessPayment')

need = ['consumer_key', 'consumer_secret', 'b2c_shortcode', 'initiator_name', 'initiator_password']
missing = [k for k in need if not vals.get(k)]
if missing:
    print('Missing from creds:', ', '.join(missing), file=sys.stderr)
    sys.exit(1)

secret_keys = {'consumer_key', 'consumer_secret', 'initiator_password'}
sql_lines = []
for key, val in vals.items():
    is_secret = 1 if key in secret_keys else 0
    esc = val.replace('\\', '\\\\').replace("'", "''")
    sql_lines.append(
        f"INSERT INTO disbursement_channel_config "
        f"(disbursement_channel_id, config_key, config_value, is_secret, created_at, updated_at, created_by, updated_by) "
        f"VALUES ({channel_id}, '{key}', '{esc}', {is_secret}, NOW(3), NOW(3), 1, 1) "
        f"ON DUPLICATE KEY UPDATE config_value=VALUES(config_value), is_secret=VALUES(is_secret), updated_at=NOW(3), updated_by=1;"
    )
sql = '\n'.join(sql_lines)
dot = load_dotenv(env_path)
user = dot.get('DB_USER', 'edairy')
password = dot.get('DB_PASSWORD', '')
db = dot.get('DB_NAME', 'dev_dairy_erp')
host = dot.get('DB_HOST', '127.0.0.1')
if host == 'host.docker.internal':
    host = '127.0.0.1'
cmd = ['mysql', '-h', host, '-u', user, f'-p{password}', db]
subprocess.run(cmd, input=sql.encode(), check=True)
print(f'Seeded {len(sql_lines)} MOBILE_MPESA config keys for channel {channel_id}')
PY
