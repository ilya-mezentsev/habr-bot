#!/usr/bin/env bash

function prepareFolders() {
  mkdir -p "$1"/test_report
  mkdir -p "$1"/data
}

function prepareFiles() {
  rm "$1"/.env 2>/dev/null
  touch "$1"/.env

  rm "$1"/data/data.db 2>/dev/null
  touch "$1"/data/data.db
}

function installGolangDeps() {
  export GOPATH="$1"
  for package in \
    github.com/jmoiron/sqlx \
    github.com/mattn/go-sqlite3 \
    golang.org/x/net/html \
    github.com/go-telegram-bot-api/telegram-bot-api
  do
    go get -v $package
  done
}

rootFolder="$1"
if [[ ${rootFolder} = '' ]]; then
  echo 'root folder was not provided'
  echo 'usage bash prepare_workspace.sh ROOT_FOLDER'
  exit 1
fi

declare -A env=(
  ['ENV_VARS_WERE_SET']="1"
  ['PROJECT_ROOT']="${rootFolder}"
  ['REPORT_FOLDER']="${rootFolder}/test_report"
  ['DB_FILE']="${rootFolder}/data/data.db"
  ['GOPATH']="${rootFolder}"
  ['TG_BOT_TOKEN']=''
  ['CATEGORIES']='go,infosecurity,programming,webdev,javascript,python'
  ['ARTICLES_RESOURCE']='https://habr.com/ru/hub'
  ['ARTICLES_FILTER']='top/monthly'
  ['ARTICLE_LINK_CLASS_NAME']='post__title_link'
)

prepareFolders "${rootFolder}"
prepareFiles "${rootFolder}"
installGolangDeps "${rootFolder}"

for envKey in "${!env[@]}"
do
  echo "${envKey}=${env[$envKey]}" >> "${rootFolder}"/.env
done
