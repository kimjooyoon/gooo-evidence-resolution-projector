#!/usr/bin/env bash
set -euo pipefail

input=${1:?go test JSON input is required}
summary=${2:?summary output is required}
metrics=${3:?metrics output is required}
mkdir -p "$(dirname "$summary")"

if [[ ! -s "$input" ]]; then
  printf '%s\n' '{"total":0,"selected":0,"executed":0,"reused":0,"failed":0,"unknown":1}' > "$summary"
else
  jq -s '
    (map(select(.Action == "run" and .Test != null)) | unique_by([.Package, .Test])) as $selected |
    (map(select(.Action == "pass" and .Test != null and .Cached == true)) | unique_by([.Package, .Test])) as $reused |
    (map(select(.Action == "fail" and .Test != null)) | unique_by([.Package, .Test])) as $failed |
    (map(select(.Action == "fail" and .Test == null)) | length) as $package_failures |
    {total: ($selected | length), selected: ($selected | length), executed: (($selected | length) - ($reused | length)), reused: ($reused | length), failed: ($failed | length), unknown: $package_failures}
  ' "$input" > "$summary"
fi

tmp_file="${metrics}.tmp"
jq --slurpfile tests "$summary" '.tests = $tests[0]' "$metrics" > "$tmp_file"
mv "$tmp_file" "$metrics"

