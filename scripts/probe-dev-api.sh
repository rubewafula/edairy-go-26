#!/usr/bin/env bash
# Probe authenticated API endpoints and report failures.
# Usage: ./probe-dev-api.sh <base_url> <email> <password>
set -euo pipefail

BASE="${1:?base URL e.g. https://api.dev.edairy.africa}"
EMAIL="${2:?email}"
PASS="${3:?password}"

login=$(curl -sf -X POST "$BASE/api/login" \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASS\"}") || { echo "LOGIN FAILED"; exit 1; }

TOKEN=$(echo "$login" | python3 -c "import sys,json; print(json.load(sys.stdin).get('token',''))")
if [[ -z "$TOKEN" ]]; then echo "No token in login response: $login"; exit 1; fi

endpoints=(
  "/api/admin-dashboard"
  "/api/member-dashboard"
  "/api/produce-dashboard"
  "/api/daily-milk-variances?page=1&limit=10"
  "/api/milk-journals?page=1&limit=10"
  "/api/milk-journal-entries?page=1&limit=10"
  "/api/milk-deliveries?page=1&limit=10"
  "/api/milk-rejects?page=1&limit=10"
  "/api/milk-local-sales?page=1&limit=10"
  "/api/cooler-milk-collections?page=1&limit=10"
  "/api/members?page=1&limit=5"
  "/api/employees?page=1&limit=5"
  "/api/finance-dashboard"
  "/api/hr-dashboard"
  "/api/payroll-dashboard"
)

fail=0
for ep in "${endpoints[@]}"; do
  code=$(curl -s -o /tmp/probe_body.txt -w '%{http_code}' \
    -H "Authorization: Bearer $TOKEN" "$BASE$ep")
  if [[ "$code" =~ ^2 ]]; then
    echo "OK  $code $ep"
  else
    echo "ERR $code $ep"
    head -c 200 /tmp/probe_body.txt; echo
    fail=$((fail + 1))
  fi
done

echo "Done: $fail failures"
exit "$fail"
