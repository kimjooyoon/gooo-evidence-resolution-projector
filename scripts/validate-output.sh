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
  reviewer-view.md
  language-maintainer-view.md
  projection-dossier.md
  projection-artifact.json
)

for name in "${expected[@]}"; do
  test -f "$output_dir/$name"
done
actual_count=$(find "$output_dir" -maxdepth 1 -type f -print | wc -l | tr -d ' ')
test "$actual_count" -eq 7

jq -e . "$output_dir/projection-manifest.json" >/dev/null
jq -e . "$output_dir/evidence-graph.json" >/dev/null
jq -e . "$output_dir/projection-receipt.json" >/dev/null
jq -e . "$output_dir/projection-artifact.json" >/dev/null

graph_cases=$(jq '.cases | length' "$output_dir/evidence-graph.json")
test "$graph_cases" -eq 9
test "$(jq '[.cases[] | select(.decision == "CLOSED")] | length' "$output_dir/evidence-graph.json")" -eq 3
test "$(jq '[.cases[] | select(.decision == "UNKNOWN")] | length' "$output_dir/evidence-graph.json")" -eq 3
test "$(jq '[.cases[] | select(.decision == "REFUTED")] | length' "$output_dir/evidence-graph.json")" -eq 3
test "$(wc -l < "$output_dir/projection-events.ndjson" | tr -d ' ')" -eq 36

test "$(sha256sum "$output_dir/evidence-graph.json" | awk '{print $1}')" = "$(jq -r '.canonical_graph_sha256' "$output_dir/projection-manifest.json")"
test "$(jq -r '.decision_digest' "$output_dir/evidence-graph.json")" = "$(jq -r '.decision_digest' "$output_dir/projection-manifest.json")"
test "$(jq -r '.inventory.caller_owned_output_kinds' "$output_dir/projection-manifest.json")" -eq 11
test "$(jq -r '.inventory.input_repo_mutations' "$output_dir/projection-manifest.json")" -eq 0
test "$(jq -r '.inventory.source_mutations' "$output_dir/projection-manifest.json")" -eq 0
test "$(jq -r '.inventory.runtime_side_effects' "$output_dir/projection-manifest.json")" -eq 0

jq -s -e '
  (map(.decision_digest) | unique | length == 1) and
  (map(.canonical_graph_sha256) | unique | length == 1) and
  (map(.omitted_field_count == (.omitted_field_ids | length)) | all) and
  (map(.visible_nodes >= 0 and .hidden_nodes >= 0 and .folded_edges >= 0 and .lost_fields >= 0) | all) and
  (group_by(.case_id) | map(map(.decision) | unique | length == 1) | all) and
  (map(select(.decision == "REFUTED") | (.fields.counterexample_ids | length > 0)) | all) and
  (map(select(.fields | has("immutable_identity") and has("authority_boundary") and has("claim_edge_ids") and has("evidence_edge_ids") and has("counterexample_edge_ids") and has("provenance_edge_ids")) | .fields.immutable_identity != null) | all) and
  (map(select(.decision == "UNKNOWN") | (.fields | has("unknown.stage") and has("unknown.step") and has("unknown.reason") and has("unknown.unknown_class") and has("unknown.next_operation") and has("unknown.blocked_by"))) | all)
' "$output_dir/projection-events.ndjson" >/dev/null

jq -e '
  .semantic_ir_schema == "gooo/evidence-resolution-projector/semantic-ir/v1" and
  (.canonical_graph_sha256 | type == "string" and length > 0) and
  .semantic_ir.schema_version == .semantic_ir_schema and
  .semantic_ir.canonical_node_count == (.semantic_ir.nodes | length) and
  .semantic_ir.canonical_edge_count == (.semantic_ir.edges | length)
' "$output_dir/projection-artifact.json" >/dev/null
jq -e '
  (.projections | length == 4) and
  (.projections | map(.role) | sort == ["AUDITOR", "LANGUAGE_MAINTAINER", "REVIEWER", "USER"]) and
  (.expansion_evaluations | map(.decision) | sort == ["CLOSED", "REFUTED", "UNKNOWN"]) and
  (.proof_cells | map(.proof_choice) | map(select(. == "FOUNDATION")) | length == 3) and
  (.proof_cells | map(.proof_choice) | map(select(. == "COHERENCE")) | length == 3) and
  (.proof_cells | map(.proof_choice) | map(select(. == "REGRESSION")) | length == 3) and
  (.proof_cells | map(.indicator) | map(select(. == "DRIVER")) | length == 3) and
  (.proof_cells | map(.indicator) | map(select(. == "OUTCOME")) | length == 3) and
  (.proof_cells | map(.indicator) | map(select(. == "GUARDRAIL")) | length == 3) and
  (.runtime.repository_writes == 0 and .runtime.source_mutations == 0 and .runtime.cross_project_required_gates == 0)
' "$output_dir/projection-artifact.json" >/dev/null
