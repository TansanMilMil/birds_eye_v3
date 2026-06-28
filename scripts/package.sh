#!/bin/bash -eu

cd "$(dirname "$0")/.."

echo "==== Create deploy package ===="
rm -rf ./backend/deploy-temp
mkdir -p ./backend/deploy-temp
mkdir -p ./backend/deploy-temp/backend/go
rsync -a --ignore-missing-args ./backend/go/dist ./backend/deploy-temp/backend/go/
rsync -a --ignore-missing-args ./backend/Dockerfile.prod ./backend/deploy-temp/backend/
rsync -a --ignore-missing-args ./docker-compose.prod.yml ./backend/deploy-temp/docker-compose.yml
rsync -a --ignore-missing-args ./frontend/dist ./backend/deploy-temp/frontend/
rsync -a --ignore-missing-args ./.env ./backend/deploy-temp/ || true
rsync -a --ignore-missing-args ./backend/init_db.sh ./backend/deploy-temp/
rsync -a --ignore-missing-args ./nginx ./backend/deploy-temp/
tar czf ./backend/deploy.tgz -C ./backend/deploy-temp .
rm -rf ./backend/deploy-temp
echo "Package created deploy.tgz"
