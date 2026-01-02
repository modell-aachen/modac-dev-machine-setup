#!/usr/bin/env bash
set -e

base_dir="$1"
claude_md_directory="$HOME/.claude"
mkdir -p "$claude_md_directory"
cp "$base_dir/templates/team-claude.md" "$claude_md_directory/CLAUDE.md"
