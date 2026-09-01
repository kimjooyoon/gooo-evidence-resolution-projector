#!/usr/bin/env bash
set -euo pipefail

output_dir=${1:?output directory is required}
expected=(
  projection-manifest.json
  evidence-graph.json
  projection-events.ndjson
  user-view.md
  operator-view.md
  auditor-view.md
  projection-receipt.json
)

for name in "${expected[@]}"; do
  test -f "$output_dir/$name"
done
actual_count=$(find "$output_dir" -maxdepth 1 -type f -print | wc -l | tr -d ' ')
test "$actual_count" -eq 7

jq -e . "$output_dir/projection-manifest.json" >/dev/null
jq -e . "$output_dir/evidence-graph.json" >/dev/null
jq -e . "$output_dir/projection-receipt.json" >/dev/null

graph_cases=$(jq '.cases | length' "$output_dir/evidence-graph.json")
test "$graph_cases" -eq 9
test "$(jq '[.cases[] | select(.decision == "CLOSED")] | length' "$output_dir/evidence-graph.json")" -eq 3
test "$(jq '[.cases[] | select(.decision == "UNKNOWN")] | length' "$output_dir/evidence-graph.json")" -eq 3
test "$(jq '[.cases[] | select(.decision == "REFUTED")] | length' "$output_dir/evidence-graph.json")" -eq 3
test "$(wc -l < "$output_dir/projection-events.ndjson" | tr -d ' ')" -eq 27

test "$(sha256sum "$output_dir/evidence-graph.json" | awk '{print $1}')" = "$(jq -r '.canonical_graph_sha256' "$output_dir/projection-manifest.json")"
test "$(jq -r '.decision_digest' "$output_dir/evidence-graph.json")" = "$(jq -r '.decision_digest' "$output_dir/projection-manifest.json")"
test "$(jq -r '.inventory.caller_owned_output_kinds' "$output_dir/projection-manifest.json")" -eq 7
test "$(jq -r '.inventory.input_repo_mutations' "$output_dir/projection-manifest.json")" -eq 0
test "$(jq -r '.inventory.source_mutations' "$output_dir/projection-manifest.json")" -eq 0
test "$(jq -r '.inventory.runtime_side_effects' "$output_dir/projection-manifest.json")" -eq 0

jq -s -e '
  (map(.decision_digest) | unique | length == 1) and
  (map(.omitted_field_count == (.omitted_field_ids | length)) | all) and
  (group_by(.case_id) | map(map(.decision) | unique | length == 1) | all) and
  (map(select(.decision == "REFUTED") | (.fields.counterexample_ids | length > 0)) | all) and
  (map(select(.decision == "UNKNOWN") | (.fields | has("unknown.stage") and has("unknown.step") and has("unknown.reason") and has("unknown.unknown_class") and has("unknown.next_operation") and has("unknown.blocked_by"))) | all)
' "$output_dir/projection-events.ndjson" >/dev/null

