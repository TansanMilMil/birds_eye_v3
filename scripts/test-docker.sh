#!/bin/bash -eu

cd "$(dirname "$0")/.."

if ! docker info >/dev/null 2>&1; then
    echo "==== docker が利用できないためローカルの go で実行します ===="
    echo "==== Run all tests ===="
    ./backend/test.sh
    exit 0
fi

echo "==== Docker compose setup ===="
docker compose down
docker compose up -d --build go
echo "==== Run all tests ===="
docker compose exec go ./test.sh
docker compose down
