#!/bin/bash -eu

echo "load .env ----------------------------"
if [ ! -f .env ]; then
    echo ".env file not found!"
else 
    source .env
fi

echo "----------------------------------------"
if [[ "${1:-}" == "--prod" ]]; then
    ssh $VENUS_SSH_HOST curl -v -X POST localhost:1111/news/scrape | jq .
else
    curl -v -X POST localhost/api/news/scrape | jq .
fi
