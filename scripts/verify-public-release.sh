#!/usr/bin/env bash
set -euo pipefail

repo=${1:?owner/repository is required}
tag=${2:?tag is required}
expected_commit=${3:?expected commit is required}
release_json=${4:?public release JSON is required}

jq -e --arg tag "$tag" '.draft == false and .prerelease == false and .tag_name == $tag and .immutable == true' "$release_json" >/dev/null
tag_object=$(gh api "repos/$repo/git/ref/tags/$tag" --jq '.object.sha')
tag_type=$(gh api "repos/$repo/git/ref/tags/$tag" --jq '.object.type')
if [[ "$tag_type" == "tag" ]]; then
  tag_commit=$(gh api "repos/$repo/git/tags/$tag_object" --jq '.object.sha')
else
  tag_commit=$tag_object
fi
test "$tag_commit" = "$expected_commit"

