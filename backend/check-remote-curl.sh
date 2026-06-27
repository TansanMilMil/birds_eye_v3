#!/bin/bash -eu

echo "load .env ----------------------------"
if [ ! -f .env ]; then
    echo ".env file not found!"
else 
    source .env
fi

ssh $VENUS_SSH_HOST curl -v localhost:1111/HealthCheck

sleep 5

echo "----------------------------------------"
target_date=$(date +%Y-%m-%d)
ssh $VENUS_SSH_HOST curl -v localhost:1111/news/${target_date} | jq .
