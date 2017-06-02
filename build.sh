#!/bin/bash

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
OUT_DIR="$SCRIPT_DIR/bin"

if [ ! -d "$OUT_DIR" ]; then
    mkdir "$OUT_DIR"
fi

cp "$SCRIPT_DIR/config.json" "$OUT_DIR/config.json"
cp -R "$SCRIPT_DIR/public" "$OUT_DIR"

go build -i -race -o "$OUT_DIR/gotunes" github.com/pulse0ne/gotunes
