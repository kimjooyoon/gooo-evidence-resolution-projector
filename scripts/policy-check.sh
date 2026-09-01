#!/usr/bin/env bash
set -euo pipefail

if command -v rg >/dev/null 2>&1; then
  search_tool=(rg -n)
else
  search_tool=(grep -R -n -E)
fi

if "${search_tool[@]}" 'secrets\.|GH_PAT|PAT_TOKEN|immutable-releases|/admin/' .github/workflows; then
  echo 'workflow policy violation: user secrets or admin/immutable API detected' >&2
  exit 1
fi
"${search_tool[@]}" '\$\{\{ github\.token \}\}' .github/workflows/release.yml >/dev/null
