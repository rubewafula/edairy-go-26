#!/bin/bash
# Live loan disbursement via M-Pesa B2C on dev (requires configure-dev-mpesa-from-creds.sh).
set -euo pipefail

API="${API:-https://api.dev.edairy.africa/api}"
EMAIL="${EMAIL:-rubewafula@gmail.com}"
PASS="${PASS:-000}"
PHONE="${PHONE:-0726986944}"
AMOUNT="${AMOUNT:-50}"
CONTRACT_ID="${CONTRACT_ID:-}"
MEMBER_ID="${MEMBER_ID:-926}"
POST_TO_GL="${POST_TO_GL:-false}"
POLL_SECS="${POLL_SECS:-120}"

JWT=$(curl -sf -X POST "$API/login" -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASS\"}" | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])")
auth=(-H "Authorization: Bearer $JWT" -H "Content-Type: application/json")

if [[ -z "$CONTRACT_ID" ]]; then
  CONTRACT_ID=$(curl -sf "$API/loan-module/contracts?status=PENDING_DISBURSEMENT&limit=5" "${auth[@]}" | python3 -c "
import sys,json
d=json.load(sys.stdin).get('data') or []
if not d: raise SystemExit('no PENDING_DISBURSEMENT contract')
print(d[0]['id'])
")
fi

REF="DEV-MPESA-LIVE-$(date +%Y%m%d%H%M%S)"
echo "Contract $CONTRACT_ID | $AMOUNT KES to $PHONE | ref $REF"

DISB=$(curl -sf -X POST "$API/loan-module/contracts/$CONTRACT_ID/disburse" "${auth[@]}" \
  -d "{
    \"amount\": $AMOUNT,
    \"channel_code\": \"MOBILE_MPESA\",
    \"channel_payload\": {\"phone_number\": \"$PHONE\"},
    \"reference\": \"$REF\",
    \"idempotency_key\": \"$REF\",
    \"post_to_gl\": $POST_TO_GL,
    \"create_milk_deduction\": false
  }")
DISB_ID=$(echo "$DISB" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])")
STATUS=$(echo "$DISB" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['status'])")
echo "Disbursement id=$DISB_ID status=$STATUS"

deadline=$((SECONDS + POLL_SECS))
while [[ $SECONDS -lt $deadline ]]; do
  sleep 5
  ROW=$(curl -sf "$API/loan-module/disbursements/$DISB_ID" "${auth[0]}")
  STATUS=$(echo "$ROW" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['status'])")
  EXT=$(echo "$ROW" | python3 -c "import sys,json; d=json.load(sys.stdin)['data']; print(d.get('external_reference') or '')")
  echo "  poll: $STATUS ext_ref=${EXT:-—}"
  if [[ "$STATUS" == "SUCCESS" || "$STATUS" == "FAILED" ]]; then
    break
  fi
done

if [[ "$STATUS" != "SUCCESS" ]]; then
  echo "Disbursement did not succeed (status=$STATUS). Check backend logs and Safaricom callback URL." >&2
  exit 1
fi

MSG="Dear member, your eDairy loan disbursement of KES ${AMOUNT} via M-Pesa is complete. Ref: ${REF}. Thank you."
curl -sf -X POST "$API/sms-outboxes" "${auth[@]}" \
  -d "$(python3 -c "import json; print(json.dumps({'sms_campaign_id':int('${SMS_CAMPAIGN_ID:-14}'),'phone_no':'$PHONE','message':'''$MSG'''}))")" \
  | python3 -c "import sys,json; m=json.load(sys.stdin); print('SMS outbox:', m.get('data', m))"

echo "Done: loan disbursement SUCCESS + notification SMS submitted (sender: edairy provider on dev)."
