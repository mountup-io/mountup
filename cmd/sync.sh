#!/usr/bin/env bash
set -euo pipefail

# Args:
# $1=hostname
# $2=keypath

export WORK_DIR="$HOME/mountup"

fswatch -r "${WORK_DIR}" | while read f; do
    rsync "${WORK_DIR}" -azP -e "ssh -i $2" "$1":~
done
