#!/usr/bin/env bash
set -e

base_path="$(cd "$(dirname "$0")" && pwd)/.."

claude_md_directory="$HOME/.claude"
mkdir -p "$claude_md_directory"
if [[ ! -f "$claude_md_directory/CLAUDE.md" ]]; then
    cp "$base_path/templates/team-claude.md" "$claude_md_directory/CLAUDE.md"
    chmod 664 "$claude_md_directory/CLAUDE.md"
fi
