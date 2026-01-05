#!/usr/bin/env bash
set -e

base_path="$(cd "$(dirname "$0")" && pwd)/.."

claude_md_directory="$HOME/.claude"
mkdir -p "$claude_md_directory"
cp "$base_path/templates/team-claude.md" "$claude_md_directory/CLAUDE.md"
