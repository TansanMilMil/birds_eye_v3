#!/bin/bash -eu

cd "$(dirname "$0")/.."
source .env

TARGET_DIR="${VENUS_HOME}/birds_eye_v3"

echo "==== Transfer to remote ===="
if [ -z "${VENUS_SSH_HOST:-}" ]; then
  echo "Error: VENUS_SSH_HOST is not set. Please check .env file."
  exit 1
fi
if [ -z "${VENUS_HOME:-}" ]; then
  echo "Error: VENUS_HOME is not set. Please check .env file."
  exit 1
fi
ssh "${VENUS_SSH_HOST}" mkdir -p "${TARGET_DIR}"
scp ./backend/deploy.tgz "${VENUS_SSH_HOST}:${TARGET_DIR}/deploy.tgz"
ssh "${VENUS_SSH_HOST}" "tar xzf ${TARGET_DIR}/deploy.tgz -C ${TARGET_DIR} && rm ${TARGET_DIR}/deploy.tgz"
rm ./backend/deploy.tgz
echo ""
echo "Current files on remote:"
ssh "${VENUS_SSH_HOST}" ls -alh "${TARGET_DIR}"
