#!/usr/bin/env bash

cd "${PROJECT_ROOT}" || exit

docker run -d -v "$(pwd)"/data:/var/www/apps/habr-bot/data --env-file .env -e DB_FILE=/var/www/apps/habr-bot/data/data.db habr-bot:latest
