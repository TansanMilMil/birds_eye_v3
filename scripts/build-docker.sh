#!/bin/bash -eu

cd "$(dirname "$0")/.."

echo "==== Docker compose setup ===="
docker compose down
docker compose up -d go
echo "==== Build ===="
docker compose exec go ./build.sh
echo "==== Docker compose teardown ===="
docker compose down
