#!/bin/bash

set -euo pipefail

FILE="$1"

# clear out old files
rm -f "$FILE-encoded.huff" "$FILE-decoded.txt"

go run main.go encode -i "$FILE.txt" -o "$FILE-encoded.huff"
go run main.go decode -i "$FILE-encoded.huff" -o "$FILE-decoded.txt"

diff -u "$FILE.txt" "$FILE-decoded.txt" | delta
