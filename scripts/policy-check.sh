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

jq -e '
  .append_only == true and
  .previous_released_contract.case_denominator == 9 and
  .current_contract.case_denominator_minimum >= .previous_released_contract.case_denominator and
  (.current_contract.reader_roles_minimum | sort) == ["AUDITOR", "LANGUAGE_MAINTAINER", "REVIEWER", "USER"] and
  .invariants.runtime.repository_writes == 0 and
  .invariants.runtime.source_mutations == 0 and
  .invariants.runtime.cross_project_required_gates == 0
' contracts/denominator-v1.json >/dev/null
