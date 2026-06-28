#!/bin/bash -eu

cd "$(dirname "$0")/.."

echo "==== Docker compose setup ===="
docker compose down
docker compose up -d go
echo "==== Run all tests ===="
docker compose exec go ./test.sh
docker compose down
