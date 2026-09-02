#!/usr/bin/env bash
set -euo pipefail

release_json=${1:?release JSON is required}
output_dir=${2:?output directory is required}

expected=(
  projection-manifest.json
  evidence-graph.json
  projection-events.ndjson
  user-view.md
  operator-view.md
  auditor-view.md
  projection-receipt.json
  reviewer-view.md
  language-maintainer-view.md
  projection-dossier.md
  projection-artifact.json
)
test "$(jq '.assets | length' "$release_json")" -eq 7
for name in "${expected[@]}"; do
  local_path="$output_dir/$name"
  local_size=$(wc -c < "$local_path" | tr -d ' ')
  local_digest=$(sha256sum "$local_path" | awk '{print $1}')
  remote_size=$(jq -r --arg name "$name" '.assets[] | select(.name == $name) | .size' "$release_json")
  remote_digest=$(jq -r --arg name "$name" '.assets[] | select(.name == $name) | .digest' "$release_json")
  test "$remote_size" = "$local_size"
  test "$remote_digest" = "sha256:$local_digest"
done
