#!/bin/sh
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
	go run main.go completion "$sh" >"completions/files-cli.$sh"
done

# Check if the generated files exist before proceeding
for sh in bash zsh fish; do
	if [ ! -f "completions/files-cli.$sh" ]; then
		echo "Error: completions/files-cli.$sh does not exist."
		exit 1
	fi
done