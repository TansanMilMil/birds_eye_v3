#!/bin/bash -eu

cd `dirname $0`
if [ -f .env ]; then
    source .env
else
    echo ".env file not found."
fi
docker compose up -d

echo 'init completed!'
