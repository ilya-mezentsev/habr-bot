ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SOURCE_PATH := $(ROOT_DIR)/source
LIBS_PATH := $(ROOT_DIR)/libs
ENV_FILE := $(ROOT_DIR)/.env
DATA_DIR := $(ROOT_DIR)/data
REPORT_FOLDER := $(ROOT_DIR)/test_report

DB_FILE_PATH := $(DATA_DIR)/data.db
ENTRYPOINT_FILE := $(ROOT_DIR)/main

ifneq (,$(wildcard ./.env))
    include .env
    export
    ENV_FILE_PARAM = --env-file .env
endif

workspace: env install-libs

build:
	unset GOPATH && cd $(ROOT_DIR) && go build main.go

test:
	unset GOPATH && cd $(SOURCE_PATH) && go test ./... -cover | { grep -v "no test files"; true; }

run-tg:
	$(ENTRYPOINT_FILE) -mode tg

run-cli:
	$(ENTRYPOINT_FILE) -mode cli

check:
	unset GOPATH && cd $(SOURCE_PATH) && go vet ./...

install-libs:
	unset GOPATH && GOMODCACHE=$(LIBS_PATH) go mod download

calc-lines:
	( find $(SOURCE_PATH) -name '*.go' -print0 | xargs -0 cat ) | wc -l

env: clean-env clean-db-file create-db-file
	echo "CATEGORIES=go,rust,infosecurity,programming,webdev,javascript,python" >> $(ENV_FILE)
	echo "ARTICLE_LINK_CLASS_NAME=post__title_link" >> $(ENV_FILE)
	echo "DB_FILE=$(DB_FILE_PATH)" >> $(ENV_FILE)
	echo "TG_BOT_TOKEN=" >> $(ENV_FILE)
	echo "ARTICLES_RESOURCE=https://habr.com/ru/hub" >> $(ENV_FILE)
	echo "CATEGORIES_FILTERS=top10,top,top/monthly" >> $(ENV_FILE)
	echo "REPORT_FOLDER=$(REPORT_FOLDER)" >> $(ENV_FILE)

create-db-file: clean-db-file
	mkdir -p $(DATA_DIR) && touch $(DB_FILE_PATH)

clean-db-file:
	rm -f $(DB_FILE_PATH)

clean-env:
	rm -f $(ENV_FILE)
