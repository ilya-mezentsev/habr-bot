ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SOURCE_PATH := $(ROOT_DIR)/source
LIBS_PATH := $(ROOT_DIR)/libs
ENV_FILE := $(ROOT_DIR)/.env
DATA_DIR := $(ROOT_DIR)/data

CONFIGS_FILE_PATH ?= $(ROOT_DIR)/config/main.json
DB_FILE_PATH := $(DATA_DIR)/data.db
ENTRYPOINT_FILE ?= $(ROOT_DIR)/main
TG_BOT_TOKEN ?= ""

ifneq (,$(wildcard ./.env))
	include .env
	export
	ENV_FILE_PARAM = --env-file .env
endif

workspace: env install-libs

build:
	unset GOPATH && cd $(ROOT_DIR) && go build main.go

test:
	unset GOPATH && cd $(SOURCE_PATH) && go test -cover -p 1 ./... | { grep -v "no test files"; true; }

run-tg:
	$(ENTRYPOINT_FILE) -mode tg -config $(CONFIGS_FILE_PATH)

run-cli:
	$(ENTRYPOINT_FILE) -mode cli -config $(CONFIGS_FILE_PATH)

check:
	unset GOPATH && cd $(SOURCE_PATH) && go vet ./...

install-libs:
	unset GOPATH && GOMODCACHE=$(LIBS_PATH) go mod download

calc-lines:
	( find $(SOURCE_PATH) -name '*.go' -print0 | xargs -0 cat ) | wc -l

env: clean-env reset-db-file
	echo "DB_FILE=$(DB_FILE_PATH)" >> $(ENV_FILE)
	echo "TG_BOT_TOKEN=$(TG_BOT_TOKEN)" >> $(ENV_FILE)

reset-db-file: clean-db-file
	mkdir -p $(DATA_DIR) && touch $(DB_FILE_PATH)

clean-db-file:
	rm -f $(DB_FILE_PATH)

clean-env:
	rm -f $(ENV_FILE)
