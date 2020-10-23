#!/usr/bin/env bash

cd "${PROJECT_ROOT}" || exit

GOPATH="${PROJECT_ROOT}" go build main.go
docker build -t habr-bot .
