#!/usr/bin/env bash
set -euo pipefail

# Args:
# $1=hostname
# $2=destDir, default ~
# $3=keypath
# $4=servername
# $5=true if push, false if pull

export WORK_DIR="$HOME/mountup/$4/"

mkdir -p "$WORK_DIR"

if [ "$5" = "push" ]; then
  # push files
  echo "push"
#  rsync "${WORK_DIR}" -azP --exclude=".*" -e "ssh -i $3" "$1":"$2"
else
  # pull files
  # shellcheck disable=SC2115
  rm -rf "${WORK_DIR}/*"
  rsync -azP --exclude=".*" -e "ssh -i $3" "$1":"$2" "${WORK_DIR}"
fi

fswatch -r "${WORK_DIR}" | while read f; do
  rsync "${WORK_DIR}" -azP --exclude=".*" -e "ssh -i $3" "$1":"$2"
done
