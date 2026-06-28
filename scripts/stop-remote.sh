#!/bin/bash -eu

cd "$(dirname "$0")/.."
source .env

TARGET_DIR="${VENUS_HOME}/birds_eye_v3"

echo "==== Stop birdseyeapi on remote ===="
if [ -z "${VENUS_SSH_HOST:-}" ]; then
  echo "Error: VENUS_SSH_HOST is not set. Please check .env file."
  exit 1
fi
ssh "${VENUS_SSH_HOST}" "docker compose -f ${TARGET_DIR}/docker-compose.yml down || true"
