#!/bin/bash

BASE_DIRECTORY="/home/galt/code/brothers_in_batash"
TEMP_FILE=$(mktemp)

# Define files to copy
FILES=("internal/pkg/models/shift.go"
"internal/app/webserver/controllers/shift_controller_test.go"
"internal/app/webserver/controllers/shift_template_controller.go")

cd "$BASE_DIRECTORY" || exit 1

for file in "${FILES[@]}"; do
    if [ -f "$file" ]; then
        echo "File: $file" >> "$TEMP_FILE"
        cat "$file" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
    else
        echo "File not found: $file"
    fi
done

# Copy the contents to the clipboard
xclip -selection clipboard < "$TEMP_FILE"

rm "$TEMP_FILE"

cd - > /dev/null || exit 1