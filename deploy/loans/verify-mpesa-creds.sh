#!/bin/bash
# Quick Daraja check: OAuth + B2C v3 dry-run (does not complete a payment if Safaricom accepts).
set -euo pipefail

CREDS="${1:-/home/reuben/apps/mpesa.creds/mpesa.env}"
CERT="${2:-${MPESA_CERT:-/home/reuben/apps/mpesa.cer}}"

if [[ ! -f "$CREDS" ]]; then
  echo "Missing: $CREDS" >&2
  exit 1
fi
if [[ ! -f "$CERT" ]]; then
  echo "Missing cert: $CERT" >&2
  exit 1
fi

python3 - "$CREDS" "$CERT" <<'PY'
import sys, json, base64, uuid, urllib.request, urllib.error
from pathlib import Path
from cryptography import x509
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.asymmetric import padding

def load_env(p):
    d = {}
    for line in Path(p).read_text().splitlines():
        s = line.strip()
        if not s or s.startswith("#") or "=" not in s:
            continue
        k, v = s.split("=", 1)
        d[k.strip()] = v.strip().strip('"').strip("'")
    return d

env = load_env(sys.argv[1])
cert_path = sys.argv[2]
need = [
    "MPESA_ENV", "MPESA_CONSUMER_KEY", "MPESA_CONSUMER_SECRET",
    "MPESA_B2C_SHORTCODE", "MPESA_INITIATOR_NAME", "MPESA_INITIATOR_PASSWORD",
]
missing = [k for k in need if not env.get(k)]
if missing:
    print("FAIL: missing keys:", ", ".join(missing))
    sys.exit(1)

prod = env.get("MPESA_ENV", "production").lower() == "production"
base = "https://api.safaricom.co.ke" if prod else "https://sandbox.safaricom.co.ke"
auth = base64.b64encode(
    f"{env['MPESA_CONSUMER_KEY']}:{env['MPESA_CONSUMER_SECRET']}".encode()
).decode()
req = urllib.request.Request(
    f"{base}/oauth/v1/generate?grant_type=client_credentials",
    headers={"Authorization": f"Basic {auth}"},
)
try:
    with urllib.request.urlopen(req, timeout=30) as r:
        tok = json.loads(r.read())["access_token"]
    print(f"OK: OAuth ({'production' if prod else 'sandbox'})")
except urllib.error.HTTPError as e:
    print(f"FAIL: OAuth HTTP {e.code}")
    sys.exit(1)

cert = Path(cert_path).read_bytes()
pub = x509.load_pem_x509_certificate(cert, default_backend()).public_key()
sec = base64.b64encode(
    pub.encrypt(env["MPESA_INITIATOR_PASSWORD"].encode(), padding.PKCS1v15())
).decode()
payload = {
    "OriginatorConversationID": str(uuid.uuid4()),
    "InitiatorName": env["MPESA_INITIATOR_NAME"],
    "SecurityCredential": sec,
    "CommandID": "BusinessPayment",
    "Amount": "10",
    "PartyA": env["MPESA_B2C_SHORTCODE"],
    "PartyB": "254700000000",
    "Remarks": "eDairy credential check",
    "QueueTimeOutURL": "https://api.dev.edairy.africa/api/mpesa-b2c-timeout",
    "ResultURL": "https://api.dev.edairy.africa/api/mpesa-b2c-result",
    "Occassion": "verify",
}
body = json.dumps(payload).encode()
url = f"{base}/mpesa/b2c/v3/paymentrequest"
req = urllib.request.Request(
    url, data=body, method="POST",
    headers={"Authorization": f"Bearer {tok}", "Content-Type": "application/json"},
)
try:
    with urllib.request.urlopen(req, timeout=60) as r:
        raw = r.read().decode()
        print("OK: B2C v3 accepted (check Safaricom response):", raw[:200])
except urllib.error.HTTPError as e:
    msg = e.read().decode(errors="replace")
    if "apiproduct match" in msg:
        print("FAIL: B2C — Daraja app has OAuth but is NOT subscribed to M-Pesa B2C (v3).")
        print("     Use consumer key/secret from an app with the B2C API product enabled,")
        print("     and match MPESA_B2C_SHORTCODE + initiator to that organization.")
    elif e.code == 401:
        print("FAIL: B2C HTTP 401 — token not valid for B2C (wrong env or missing B2C product).")
    else:
        print(f"FAIL: B2C HTTP {e.code}")
    sys.exit(1)
PY

echo "Credentials look good for live loan disbursement."
