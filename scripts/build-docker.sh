#!/bin/bash -eu

cd "$(dirname "$0")/.."

if ! docker info >/dev/null 2>&1; then
    echo "==== docker が利用できないためローカルの go で実行します ===="
    echo "==== Build ===="
    ./backend/build.sh
    exit 0
fi

echo "==== Docker compose setup ===="
docker compose down
docker compose up -d go
echo "==== Build ===="
docker compose exec go ./build.sh
echo "==== Docker compose teardown ===="
docker compose down
