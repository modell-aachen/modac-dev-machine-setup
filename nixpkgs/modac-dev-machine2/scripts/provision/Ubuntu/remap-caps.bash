#!/usr/bin/env bash
set -e

if [[ "$(gsettings get org.gnome.desktop.input-sources xkb-options)" != *"caps:escape"* ]]; then
    gsettings set org.gnome.desktop.input-sources xkb-options "['caps:escape']"
fi
