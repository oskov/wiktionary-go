#/bin/bash

#!/bin/bash

# Define the URL and output file path
URL="https://en.wiktionary.org/api/rest_v1/?spec"
OUTPUT_FILE="./api/wiktionary/openapi.json"

# Create the output directory if it doesn't exist
mkdir -p "$(dirname "$OUTPUT_FILE")"

# Download the file
curl -o "$OUTPUT_FILE" "$URL"

# Check if the download was successful
if [ $? -eq 0 ]; then
    echo "File downloaded successfully to $OUTPUT_FILE"
else
    echo "Failed to download the file from $URL"
    exit 1
fi