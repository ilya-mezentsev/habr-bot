#!/usr/bin/env bash

cd "${PROJECT_ROOT}" || exit
docker build -t habr-bot .
