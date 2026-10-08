#!/bin/bash -eu

echo "load .env ----------------------------"
if [ ! -f .env ]; then
    echo ".env file not found!"
else 
    source .env
fi

echo "----------------------------------------"
if [[ "${1:-}" == "--prod" ]]; then
    ssh $VENUS_SSH_HOST docker exec birds_eye_go curl -sS -X POST localhost:8080/news/scrape | jq .
else
    docker compose exec go curl -sS -X POST localhost:8080/news/scrape | jq .
fi
