#!/bin/bash
# Script to make test files

TARGET_DIR="testfiles"

extensions=(
    "txt" "md" "log"             # Documents
    "jpg" "png" "gif" "svg"      # Pictures
    "pdf" "docx" "xlsx" "pptx"   # Office
    "mp3" "wav" "mp4" "mkv"      # Media
    "zip" "tar.gz" "rar" "7z"    # Archives
    "html" "css" "js" "json"     # Web
    "sh" "py" "cpp" "go"         # Code
    "conf" "ini" "yaml" "xml"    # Config
)

mkdir -p "$TARGET_DIR"

echo "Generating into '$TARGET_DIR'..."

count=0

for ext in "${extensions[@]}"; do
    for i in {1..3}; do
        dd if=/dev/urandom of="$TARGET_DIR/$i.$ext" bs=1k count=1 status=none
        count=$((count + 1))
    done
done

echo "Ready! $count files with ${#extensions[@]} different exts."
