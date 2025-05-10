#!/usr/bin/env bash
ui_array=("$@")

rm -rf ./test/.credentials
rm -rf .credentials

# Iterating through elements
for ui in "${ui_array[@]}"; do
    rm -f cmd/"$ui"/"$(APP_NAME)"
    app_tag="$(APP_NAME)"-"$ui":"$(VERSION)"
    docker image rm "$(docker image inspect "$app_tag" | jq '.[].Id' | rg -o "sha256:(.{12})" | cut -d: -f2)"
done
