#!/usr/bin/env bash
set -euo pipefail

step_name=${1:?step name is required}
metrics_file=${2:?metrics file is required}
shift 2

mkdir -p "$(dirname "$metrics_file")"
if [[ ! -f "$metrics_file" ]]; then
  printf '%s\n' '{"schema_version":"1.0.0","steps":{}}' > "$metrics_file"
fi

time_file=$(mktemp)
start_ms=$(date +%s%3N)
set +e
/usr/bin/time -f 'peak_rss_kib=%M' "$@" 2> >(tee "$time_file" >&2)
status=$?
set -e
end_ms=$(date +%s%3N)
wall_ms=$((end_ms - start_ms))
peak_rss_kib=$(awk -F= '/^peak_rss_kib=/{value=$2} END{if (value == "") value=0; print value}' "$time_file")
rm -f "$time_file"

tmp_file="${metrics_file}.tmp"
jq --arg name "$step_name" --argjson wall "$wall_ms" --argjson rss "$peak_rss_kib" --argjson exit_code "$status" \
  '.steps[$name] = {wall_ms: $wall, peak_rss_kib: $rss, exit_code: $exit_code}' \
  "$metrics_file" > "$tmp_file"
mv "$tmp_file" "$metrics_file"
exit "$status"

