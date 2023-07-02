#!/bin/bash
git fetch --tags --prune-tags --prune
latest_tag=$(git tag -l --sort=-v:refname | head -n 1)

new_tag=$(echo "$latest_tag" | awk -F. '{ printf("%s.%s.%s", $1, $2, $3+1) }')

git tag "${new_tag}"

git push origin "${new_tag}"
