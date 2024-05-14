#!/bin/bash

DIRECTORY="/home/galt/code/brothers_in_batash"
TEMP_FILE=$(mktemp)

# Define exclusions: List directories and files to exclude
TREE_EXCLUDE_DIRS=(".git" ".idea" "node_modules" ".venv" "test" "e2e" "docs" "ops" ".ci" ".github")
EXCLUDE_DIRS=(".git" "deploy" "cmd" ".idea" "node_modules" ".venv" "test" "e2e" "docs" "ops" ".ci" ".github")
EXCLUDE_FILES=("go.mod" "go.sum" "Dockerfile" "makefile" ".gitignore" "README.md" "*.png" "*.ico" "package-lock.json")

if [ ! -d "$DIRECTORY" ]; then
  echo "$DIRECTORY does not exist."
  exit 1
fi

# Generate the exclude pattern for tree command
EXCLUDE_PATTERN=""
for EXCLUDE_DIR in "${TREE_EXCLUDE_DIRS[@]}"; do
    EXCLUDE_PATTERN+="-I $EXCLUDE_DIR "
done

tree "$DIRECTORY" $EXCLUDE_PATTERN > "$TEMP_FILE"

cd "$DIRECTORY" || exit 1

FIND_CMD="find . -type d"
for EXCLUDE_DIR in "${EXCLUDE_DIRS[@]}"; do
    FIND_CMD+=" \( -path './$EXCLUDE_DIR' -prune \) -o"
done
FIND_CMD+=" -type f"
for EXCLUDE_FILE in "${EXCLUDE_FILES[@]}"; do
    FIND_CMD+=" ! -name '$EXCLUDE_FILE'"
done
FIND_CMD+=" -print0"

# Execute find command and write the contents of the files to the temporary file
eval "$FIND_CMD" | while IFS= read -r -d $'\0' file; do
    echo "File: $file"
    echo "$file" >> "$TEMP_FILE"
    cat "$file" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
done

# Copy the contents to the clipboard
xclip -selection clipboard < "$TEMP_FILE"

rm "$TEMP_FILE"

cd - > /dev/null || exit 1