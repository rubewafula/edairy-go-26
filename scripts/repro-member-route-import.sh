#!/usr/bin/env bash
# Reproduce duplicate routes during member import on dev.
# Usage: ./repro-member-route-import.sh [base_url] [jwt_token]
set -euo pipefail

BASE="${1:-https://api.dev.edairy.africa}"
TOKEN="${2:?JWT token required}"

ROUTE_NAME="ReproRoute-$(date +%s)"
ROWS=40
CSV="/tmp/member-route-repro-$$.csv"

header="Member No,First Name,Last Name,Other Names,ID No,Gender,Date of Birth,Date Registered,Member Type,Route,Primary Phone,Secondary Phone,Email,Bank,Bank Branch,Account No,Account Name,Number of Cows"
echo "$header" > "$CSV"
for i in $(seq 1 "$ROWS"); do
  idno="REPRO$(date +%s)${i}"
  echo "M${i},First${i},Last${i},,${idno},MALE,,2024-01-01,Ordinary,${ROUTE_NAME},254712345678,,,,,,,1" >> "$CSV"
done

echo "Uploading $ROWS members with route '$ROUTE_NAME'..."
code=$(curl -s -o /tmp/repro-import-resp.txt -w '%{http_code}' \
  -X POST "$BASE/api/members/import" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@$CSV")
echo "Import HTTP $code"
cat /tmp/repro-import-resp.txt
echo

echo "Waiting for background import..."
sleep 8

echo "ROUTE_NAME=$ROUTE_NAME"
echo "CSV=$CSV"
