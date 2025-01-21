#!/bin/bash

file_path="manifest.yaml"

# Append multiple lines to the file
cat <<EOF >> "$file_path"

DescriberTag: local-$TAG
UpdateDate: $(date +%Y-%m-%d)

EOF

echo "Data has been appended to $file_path"