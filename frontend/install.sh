#!/bin/bash

cd `dirname $0`

docker compose up -d
docker compose exec node npm install

