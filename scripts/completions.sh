#!/bin/sh
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
	go run ./cmd/glide/main.go completion "$sh" >"completions/glide.$sh"
done
