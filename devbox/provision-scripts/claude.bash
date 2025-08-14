#!/bin/bash
set -e

claude_md_directory="$HOME/.claude"
mkdir -p "$claude_md_directory"
cp "./templates/team-claude.md" "$claude_md_directory/CLAUDE.md"
