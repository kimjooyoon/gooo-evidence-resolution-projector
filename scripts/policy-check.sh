#!/usr/bin/env bash
set -euo pipefail

if rg -n 'secrets\.|GH_PAT|PAT_TOKEN|immutable-releases|/admin/' .github/workflows; then
  echo 'workflow policy violation: user secrets or admin/immutable API detected' >&2
  exit 1
fi
rg -n '\$\{\{ github\.token \}\}' .github/workflows/release.yml >/dev/null

