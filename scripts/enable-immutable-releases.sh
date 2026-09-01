#!/usr/bin/env bash
set -euo pipefail

repo=${1:?owner/repository is required}
enabled=$(gh api "repos/$repo/immutable-releases" --jq '.enabled' 2>/dev/null || true)
if [[ "$enabled" == "true" ]]; then
  echo "immutable releases already enabled for $repo"
  exit 0
fi
gh api --method PUT \
  -H 'Accept: application/vnd.github+json' \
  -H 'X-GitHub-Api-Version: 2026-03-10' \
  "repos/$repo/immutable-releases" >/dev/null
echo "immutable releases enabled for $repo"
