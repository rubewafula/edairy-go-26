#!/usr/bin/env bash
# Migrate all six production dairies from systemd to Docker Compose.
# Usage: sudo EDAIRY_VERSION=1.0.0 ./migrate-all-to-docker.sh
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EDAIRY_VERSION="${EDAIRY_VERSION:-1.0.0}"
export EDAIRY_VERSION

declare -a DAIRIES=(
  "nguene 8510 8511"
  "mutuati 8508 8509"
  "mwimbi 8506 8507"
  "mukululu 8504 8505"
  "tigania-west 8503 8502"
  "arithi 8500 8501"
)

echo "Migrating ${#DAIRIES[@]} dairies to Docker (EDAIRY_VERSION=${EDAIRY_VERSION})"
echo "Order: nguene first (rehearsal), then remaining five."

for entry in "${DAIRIES[@]}"; do
  read -r slug api_port ui_port <<< "$entry"
  echo ""
  echo "========================================"
  echo "Migrating: ${slug} (api=${api_port}, ui=${ui_port})"
  echo "========================================"
  "${SCRIPT_DIR}/create-dairy.sh" --migrate --skip-db --skip-haproxy "${slug}" "${api_port}" "${ui_port}"
done

echo ""
echo "All dairies migrated. Verify each stack:"
for entry in "${DAIRIES[@]}"; do
  read -r slug api_port ui_port <<< "$entry"
  echo "  curl http://127.0.0.1:${api_port}/health  # ${slug} API"
  echo "  curl http://127.0.0.1:${ui_port}/runtime-config.js  # ${slug} UI"
done
