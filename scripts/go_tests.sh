#!/usr/bin/env bash
if [[ ${ENV_VARS_WERE_SET} != '1' ]]; then
  echo 'env variables are not set'
  exit 1
fi
# we need to check that REPORT_FOLDER is subdirectory for GOPATH
if [[ ${REPORT_FOLDER}/ != ${GOPATH}* ]]; then
  echo 'go tests report folder should be in GOPATH'
  exit 1
fi

folders=()
cd "${GOPATH}"/src || exit
for dir in $(find . -type d)
do
  # should skip github libs
  if [[ ${dir} == *github* || ${dir} == *golang.org* ]]; then
    continue
  fi
  if tests=$(find "${GOPATH}"/src/"${dir}" -maxdepth 1 -name '*_test.go'); [[ ${tests} != "" ]]; then
    folders+=("${dir}")
  fi
done

rm -rf "${REPORT_FOLDER:?}"/*
for dir in "${folders[@]}"
do
  reportFileName=$(echo -n "${dir}" | md5sum | awk '{print $1}')
  reportFilePath=${REPORT_FOLDER}/${reportFileName}
  cd "${GOPATH}"/src/"${dir}" && go test -coverprofile="${reportFilePath}".out
  if [[ $1 = html ]]; then # open reports in browser
    go tool cover -html="${reportFilePath}".out -o "${reportFilePath}".html
    chromium "${reportFilePath}".html >/dev/null 2>&1 &
  fi
done
