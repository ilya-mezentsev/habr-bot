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

ips=(
  108.61.190.167:8080
  136.244.86.8:8080
  185.17.123.14:3128
  37.187.3.175:3128
  5.9.252.251:3128
  51.91.212.159:3128
  151.80.199.89:3128
  82.119.170.106:8080
  209.97.138.116:8080
  51.255.103.170:3129
  178.169.198.238:8080
  217.182.253.119:3130
  54.37.131.45:3128
  80.90.80.54:8080
  80.187.140.26:8080
)
joined_ips=$(printf ",%s" "${ips[@]}")
joined_ips=${joined_ips:1}

declare -A env=(
  ['ENV_VARS_WERE_SET']="1"
  ['PROJECT_ROOT']="${rootFolder}"
  ['REPORT_FOLDER']="${rootFolder}/test_report"
  ['DB_FILE']="${rootFolder}/data/data.db"
  ['GOPATH']="${rootFolder}"
  ['TG_BOT_TOKEN']=''
  ['TG_HTTP_PROXY_IPS']=$joined_ips
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
