#!/usr/bin/env bash
ui_array=("$@")

if ! $RELEASE; then
    for ui in "${ui_array[@]}"; do
        go build -o cmd/"$ui"/"$(APP_NAME)" cmd/"$ui"/main.go
    done
fi

if $RELEASE; then
    for ui in "${ui_array[@]}"; do
        docker build --tag "$(APP_NAME)"-"$ui":"$(VERSION)" --build-arg VERSION="$(VERSION)" --build-arg NAME="$(APP_NAME)" --target "$ui" ./build/
    done
fi
