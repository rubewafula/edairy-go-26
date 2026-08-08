#!/bin/bash
# UI-parity E2E: new loan for member 820 (M8) on dev
set -euo pipefail
API="${API:-http://127.0.0.1:8512/api}"
MEMBER_ID="${MEMBER_ID:-820}"
DEV_PASS="$1"

JWT=$(curl -sf -X POST "$API/login" -H "Content-Type: application/json" \
  -d '{"email":"rubewafula@gmail.com","password":"000"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])")

auth=(-H "Authorization: Bearer $JWT" -H "Content-Type: application/json")

echo "=== Member $MEMBER_ID ==="
mysql -uedairy -p"$DEV_PASS" dev_dairy_erp -N -e \
  "SELECT member_no, first_name, last_name FROM member_registrations WHERE id=$MEMBER_ID;"

echo "=== 1. Create application (UI: Loan Applications > Create) ==="
APP=$(curl -sf -X POST "$API/loan-module/applications" "${auth[@]}" \
  -d "{\"member_id\":$MEMBER_ID,\"loan_product_id\":1,\"requested_amount\":75000,\"requested_term_months\":6,\"purpose\":\"UI test - dairy feed\",\"expected_disbursement_date\":\"2026-08-15\"}")
APP_ID=$(echo "$APP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])")
echo "Application ID: $APP_ID status=$(echo "$APP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['status'])")"

echo "=== 2. Submit (UI: Application Show > Submit) ==="
curl -sf -X POST "$API/loan-module/applications/$APP_ID/submit" "${auth[@]}" | python3 -c "import sys,json; print('status:', json.load(sys.stdin)['data']['status'])"

echo "=== 3. Approve (UI: Application Show > Approve) ==="
APPROVE=$(curl -sf -X POST "$API/loan-module/applications/$APP_ID/approve" "${auth[@]}" \
  -d '{"approval_role":"CREDIT_OFFICER","approved_amount":75000,"approved_term_months":6,"final_approval":true,"comments":"UI E2E approve"}')
CONTRACT_ID=$(echo "$APPROVE" | python3 -c "import sys,json; d=json.load(sys.stdin)['data']; print(d.get('loan_contract_id') or '')")
echo "Contract ID: $CONTRACT_ID"

echo "=== 4. Disburse (UI: Contract Show > Disburse) ==="
REF="UI-E2E-DISB-$APP_ID"
curl -sf -X POST "$API/loan-module/contracts/$CONTRACT_ID/disburse" "${auth[@]}" \
  -d "{\"amount\":75000,\"method\":\"BANK\",\"reference\":\"$REF\",\"idempotency_key\":\"$REF\",\"post_to_gl\":true,\"create_milk_deduction\":true}" \
  | python3 -c "import sys,json; d=json.load(sys.stdin)['data']; print('disbursed:', d.get('amount'))"

echo "=== 5. Verify contract + schedule ==="
curl -sf "$API/loan-module/contracts/$CONTRACT_ID" "${auth[0]}" | python3 -c "
import sys,json
c=json.load(sys.stdin)['data']
print('contract:', c['contract_no'], 'status:', c['status'], 'outstanding:', c['outstanding_principal'])
"
curl -sf "$API/loan-module/contracts/$CONTRACT_ID/schedule" "${auth[0]}" | python3 -c "
import sys,json
lines=json.load(sys.stdin)['data']
print('schedule:', len(lines), 'installments, first due:', lines[0]['total_due'], lines[0]['status'])
"

echo "=== 6. Portfolio report ==="
curl -sf "$API/loan-module/reports/portfolio" "${auth[0]}" | python3 -c "
import sys,json
p=json.load(sys.stdin)['data']
print('active:', p['active_contracts'], 'outstanding:', p['total_outstanding'])
"

echo "=== DONE: UI test data ready ==="
echo "Application: $APP_ID | Contract: $CONTRACT_ID | Member: $MEMBER_ID (M8)"
